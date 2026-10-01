package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"vocabulary/models"
	"vocabulary/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrOAuthNotConfigured   = errors.New("oauth provider is not configured")
	ErrOAuthProfileNotFound = errors.New("oauth profile not found")
)

type AuthService interface {
	CreateGuestUser(ctx context.Context, displayName string) (*models.User, error)
	Register(ctx context.Context, request models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, request models.LoginRequest) (*models.AuthResponse, error)
	GoogleLoginURL() (string, error)
	FacebookLoginURL() (string, error)
	LoginWithGoogleCode(ctx context.Context, code string) (*models.AuthResponse, error)
	LoginWithFacebookCode(ctx context.Context, code string) (*models.AuthResponse, error)
}

type authService struct {
	userRepository repositories.UserRepository
	httpClient     *http.Client
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *authService) CreateGuestUser(ctx context.Context, displayName string) (*models.User, error) {
	user := &models.User{
		DisplayName: strings.TrimSpace(displayName),
		UserType:    models.UserTypeGuest,
	}
	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create guest user: %w", err)
	}

	return user, nil
}

func (s *authService) Register(ctx context.Context, request models.RegisterRequest) (*models.AuthResponse, error) {
	email := normalizeEmail(request.Email)
	existingPasswordIdentity, err := s.userRepository.FindIdentityByProviderAndEmail(ctx, models.AuthProviderPassword, email)
	if err != nil {
		return nil, fmt.Errorf("find password identity: %w", err)
	}
	if existingPasswordIdentity != nil {
		return nil, ErrDuplicateEmail
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	passwordHashString := string(passwordHash)
	identity := &models.UserAuthIdentity{
		Provider:      models.AuthProviderPassword,
		Email:         stringPointer(email),
		EmailVerified: false,
		PasswordHash:  &passwordHashString,
	}

	existingEmailIdentity, err := s.userRepository.FindIdentityByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find identity by email: %w", err)
	}
	if existingEmailIdentity != nil {
		user := existingEmailIdentity.User
		user.DisplayName = displayNameOrFallback(request.DisplayName, user.DisplayName)
		user.UserType = models.UserTypeRegistered
		if err := s.userRepository.UpdateUser(ctx, &user); err != nil {
			return nil, fmt.Errorf("update user: %w", err)
		}

		identity.UserID = user.ID
		if err := s.userRepository.CreateIdentity(ctx, identity); err != nil {
			return nil, fmt.Errorf("create password identity: %w", err)
		}

		return s.authResponse("registered successfully", user)
	}

	user := &models.User{
		DisplayName: strings.TrimSpace(request.DisplayName),
		UserType:    models.UserTypeRegistered,
	}
	if err := s.userRepository.CreateUserWithIdentity(ctx, user, identity); err != nil {
		return nil, fmt.Errorf("create registered user: %w", err)
	}

	return s.authResponse("registered successfully", *user)
}

func (s *authService) Login(ctx context.Context, request models.LoginRequest) (*models.AuthResponse, error) {
	email := normalizeEmail(request.Email)
	identity, err := s.userRepository.FindIdentityByProviderAndEmail(ctx, models.AuthProviderPassword, email)
	if err != nil {
		return nil, fmt.Errorf("find password identity: %w", err)
	}
	if identity == nil || identity.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*identity.PasswordHash), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.authResponse("logged in successfully", identity.User)
}

func (s *authService) GoogleLoginURL() (string, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || redirectURL == "" {
		return "", ErrOAuthNotConfigured
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	if state := os.Getenv("OAUTH_STATE"); state != "" {
		params.Set("state", state)
	}

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode(), nil
}

func (s *authService) FacebookLoginURL() (string, error) {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	redirectURL := os.Getenv("FACEBOOK_REDIRECT_URL")
	if clientID == "" || redirectURL == "" {
		return "", ErrOAuthNotConfigured
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "email,public_profile")
	if state := os.Getenv("OAUTH_STATE"); state != "" {
		params.Set("state", state)
	}

	return "https://www.facebook.com/v19.0/dialog/oauth?" + params.Encode(), nil
}

func (s *authService) LoginWithGoogleCode(ctx context.Context, code string) (*models.AuthResponse, error) {
	token, err := s.exchangeGoogleCode(ctx, code)
	if err != nil {
		return nil, err
	}

	profile, err := s.fetchGoogleProfile(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.loginWithSocialProfile(ctx, models.AuthProviderGoogle, profile)
}

func (s *authService) LoginWithFacebookCode(ctx context.Context, code string) (*models.AuthResponse, error) {
	token, err := s.exchangeFacebookCode(ctx, code)
	if err != nil {
		return nil, err
	}

	profile, err := s.fetchFacebookProfile(ctx, token)
	if err != nil {
		return nil, err
	}

	return s.loginWithSocialProfile(ctx, models.AuthProviderFacebook, profile)
}

func (s *authService) loginWithSocialProfile(ctx context.Context, provider string, profile socialProfile) (*models.AuthResponse, error) {
	if profile.ID == "" {
		return nil, ErrOAuthProfileNotFound
	}

	identity, err := s.userRepository.FindIdentityByProviderAndProviderUserID(ctx, provider, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("find social identity: %w", err)
	}
	if identity != nil {
		return s.authResponse("logged in successfully", identity.User)
	}

	newIdentity := &models.UserAuthIdentity{
		Provider:       provider,
		ProviderUserID: stringPointer(profile.ID),
		Email:          optionalStringPointer(normalizeEmail(profile.Email)),
		EmailVerified:  profile.EmailVerified,
	}

	if profile.Email != "" {
		existingEmailIdentity, err := s.userRepository.FindIdentityByEmail(ctx, normalizeEmail(profile.Email))
		if err != nil {
			return nil, fmt.Errorf("find identity by email: %w", err)
		}
		if existingEmailIdentity != nil {
			user := existingEmailIdentity.User
			user.UserType = models.UserTypeRegistered
			if user.AvatarURL == nil {
				user.AvatarURL = optionalStringPointer(profile.AvatarURL)
			}
			if strings.TrimSpace(user.DisplayName) == "" {
				user.DisplayName = displayNameOrFallback(profile.Name, "Player")
			}
			if err := s.userRepository.UpdateUser(ctx, &user); err != nil {
				return nil, fmt.Errorf("update user: %w", err)
			}

			newIdentity.UserID = user.ID
			if err := s.userRepository.CreateIdentity(ctx, newIdentity); err != nil {
				return nil, fmt.Errorf("create social identity: %w", err)
			}

			return s.authResponse("logged in successfully", user)
		}
	}

	user := &models.User{
		DisplayName: displayNameOrFallback(profile.Name, "Player"),
		AvatarURL:   optionalStringPointer(profile.AvatarURL),
		UserType:    models.UserTypeRegistered,
	}
	if err := s.userRepository.CreateUserWithIdentity(ctx, user, newIdentity); err != nil {
		return nil, fmt.Errorf("create social user: %w", err)
	}

	return s.authResponse("logged in successfully", *user)
}

func (s *authService) exchangeGoogleCode(ctx context.Context, code string) (string, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return "", ErrOAuthNotConfigured
	}

	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirectURL)
	form.Set("grant_type", "authorization_code")

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := s.postFormJSON(ctx, "https://oauth2.googleapis.com/token", form, &response); err != nil {
		return "", err
	}
	if response.AccessToken == "" {
		return "", ErrOAuthProfileNotFound
	}

	return response.AccessToken, nil
}

func (s *authService) fetchGoogleProfile(ctx context.Context, accessToken string) (socialProfile, error) {
	var response struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := s.getBearerJSON(ctx, "https://www.googleapis.com/oauth2/v2/userinfo", accessToken, &response); err != nil {
		return socialProfile{}, err
	}

	return socialProfile{
		ID:            response.ID,
		Email:         response.Email,
		EmailVerified: response.VerifiedEmail,
		Name:          response.Name,
		AvatarURL:     response.Picture,
	}, nil
}

func (s *authService) exchangeFacebookCode(ctx context.Context, code string) (string, error) {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	clientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")
	redirectURL := os.Getenv("FACEBOOK_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return "", ErrOAuthNotConfigured
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("redirect_uri", redirectURL)
	params.Set("code", code)

	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := s.getJSON(ctx, "https://graph.facebook.com/v19.0/oauth/access_token?"+params.Encode(), &response); err != nil {
		return "", err
	}
	if response.AccessToken == "" {
		return "", ErrOAuthProfileNotFound
	}

	return response.AccessToken, nil
}

func (s *authService) fetchFacebookProfile(ctx context.Context, accessToken string) (socialProfile, error) {
	params := url.Values{}
	params.Set("fields", "id,name,email,picture.type(large)")
	params.Set("access_token", accessToken)

	var response struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}
	if err := s.getJSON(ctx, "https://graph.facebook.com/v19.0/me?"+params.Encode(), &response); err != nil {
		return socialProfile{}, err
	}

	return socialProfile{
		ID:            response.ID,
		Email:         response.Email,
		EmailVerified: response.Email != "",
		Name:          response.Name,
		AvatarURL:     response.Picture.Data.URL,
	}, nil
}

func (s *authService) authResponse(message string, user models.User) (*models.AuthResponse, error) {
	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &models.AuthResponse{
		Message: message,
		Token:   token,
		User:    models.NewUserResponse(user),
	}, nil
}

func (s *authService) postFormJSON(ctx context.Context, endpoint string, form url.Values, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return s.doJSON(request, target)
}

func (s *authService) getBearerJSON(ctx context.Context, endpoint string, token string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)

	return s.doJSON(request, target)
}

func (s *authService) getJSON(ctx context.Context, endpoint string, target any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	return s.doJSON(request, target)
}

func (s *authService) doJSON(request *http.Request, target any) error {
	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("oauth request failed with status %d: %s", response.StatusCode, string(body))
	}

	if err := json.NewDecoder(bytes.NewReader(body)).Decode(target); err != nil {
		return err
	}

	return nil
}

type socialProfile struct {
	ID            string
	Email         string
	EmailVerified bool
	Name          string
	AvatarURL     string
}

func generateJWT(userID uint) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	payload := map[string]any{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsignedToken := encodedHeader + "." + encodedPayload

	mac := hmac.New(sha256.New, []byte(jwtSecret()))
	if _, err := mac.Write([]byte(unsignedToken)); err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return unsignedToken + "." + signature, nil
}

func ValidateJWT(token string) (uint, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(jwtSecret()))
	if _, err := mac.Write([]byte(unsignedToken)); err != nil {
		return 0, err
	}

	expectedSignature := mac.Sum(nil)
	actualSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return 0, ErrInvalidToken
	}
	if !hmac.Equal(actualSignature, expectedSignature) {
		return 0, ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	var payload struct {
		UserID uint  `json:"user_id"`
		Exp    int64 `json:"exp"`
	}
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return 0, ErrInvalidToken
	}
	if payload.UserID == 0 {
		return 0, ErrInvalidToken
	}
	if payload.Exp <= time.Now().Unix() {
		return 0, ErrExpiredToken
	}

	return payload.UserID, nil
}

func jwtSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}

	return "development-secret"
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func displayNameOrFallback(displayName string, fallback string) string {
	displayName = strings.TrimSpace(displayName)
	if displayName != "" {
		return displayName
	}

	return fallback
}

func stringPointer(value string) *string {
	return &value
}

func optionalStringPointer(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return &value
}
