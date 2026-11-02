package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"vocabulary/middlewares"
	"vocabulary/models"
	"vocabulary/services"
	"vocabulary/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) CreateGuestUser(c *gin.Context) {
	var request models.CreateGuestUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "display_name is required"})
		return
	}

	user, err := h.authService.CreateGuestUser(c.Request.Context(), request.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot create guest user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "guest user created successfully",
		"user":    models.NewUserResponse(*user),
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	switch c.ContentType() {
	case "application/json", "application/x-www-form-urlencoded", "multipart/form-data":
	default:
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "use application/json or form data"})
		return
	}
	const maxRegistrationBytes = storage.MaxAvatarBytes + (64 << 10)
	if c.Request.ContentLength > maxRegistrationBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": "registration payload is too large"})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRegistrationBytes)
	defer func() {
		if c.Request.MultipartForm != nil {
			_ = c.Request.MultipartForm.RemoveAll()
		}
	}()

	var request models.RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": "registration payload is too large"})
			return
		}
		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid registration data"})
			return
		}
		// Revalidate after normalization so surrounding spaces in names/emails work.
	}
	if err := services.ValidateRegisterRequest(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), request)
	if err != nil {
		var validationError *services.RegistrationValidationError
		if errors.As(err, &validationError) {
			c.JSON(http.StatusBadRequest, gin.H{"message": validationError.Message})
			return
		}
		if errors.Is(err, storage.ErrAvatarTooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"message": storage.ErrAvatarTooLarge.Error()})
			return
		}
		if errors.Is(err, storage.ErrInvalidAvatar) {
			c.JSON(http.StatusBadRequest, gin.H{"message": storage.ErrInvalidAvatar.Error()})
			return
		}
		if errors.Is(err, services.ErrDuplicateEmail) {
			c.JSON(http.StatusConflict, gin.H{"message": "email already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot register user"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot login"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint(middlewares.UserIDContextKey)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
		return
	}

	response, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot get current user"})
		return
	}

	response.Token = strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	loginURL, err := h.authService.GoogleLoginURL()
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": loginURL})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	if !hasValidOAuthState(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is required"})
		return
	}

	response, err := h.authService.LoginWithGoogleCode(c.Request.Context(), code)
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) FacebookLogin(c *gin.Context) {
	loginURL, err := h.authService.FacebookLoginURL()
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": loginURL})
}

func (h *AuthHandler) FacebookCallback(c *gin.Context) {
	if !hasValidOAuthState(c) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is required"})
		return
	}

	response, err := h.authService.LoginWithFacebookCode(c.Request.Context(), code)
	if err != nil {
		writeOAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func hasValidOAuthState(c *gin.Context) bool {
	expectedState := os.Getenv("OAUTH_STATE")
	if expectedState == "" {
		return true
	}

	return c.Query("state") == expectedState
}

func writeOAuthError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrOAuthNotConfigured) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "oauth provider is not configured"})
		return
	}
	if errors.Is(err, services.ErrOAuthProfileNotFound) {
		c.JSON(http.StatusBadGateway, gin.H{"message": "cannot read oauth profile"})
		return
	}

	c.JSON(http.StatusBadGateway, gin.H{"message": "oauth login failed"})
}
