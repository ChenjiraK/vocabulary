package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vocabulary/models"
	"vocabulary/repositories"
	"vocabulary/storage"

	"golang.org/x/crypto/bcrypt"
)

func registerRequest() models.RegisterRequest {
	return models.RegisterRequest{
		DisplayName: "  Jane  ", Email: "  JANE@Example.com  ",
		Password: "password123", ConfirmPassword: "password123",
	}
}

func registerAvatar(t *testing.T) *multipart.FileHeader {
	t.Helper()
	var picture bytes.Buffer
	if err := png.Encode(&picture, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("profile_picture", "../../original.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(picture.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/auth/register", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = request.MultipartForm.RemoveAll() })
	return request.MultipartForm.File["profile_picture"][0]
}

func TestRegisterPersistsProfileAndReturnsCurrentUser(t *testing.T) {
	t.Setenv("PUBLIC_ASSET_BASE_URL", "https://assets.example.com/")
	directory := t.TempDir()
	repository := &fakeUserRepository{}
	service := NewAuthService(repository).(*authService)
	service.avatarStore = &storage.LocalAvatarStore{Directory: directory}
	request := registerRequest()
	request.ProfilePicture = registerAvatar(t)
	response, err := service.Register(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if len(repository.createdUsersWithIdentity) != 1 {
		t.Fatal("expected one user and password identity")
	}
	created := repository.createdUsersWithIdentity[0]
	if created.user.DisplayName != "Jane" || created.user.UserType != models.UserTypeRegistered || created.identity.Provider != models.AuthProviderPassword || *created.identity.Email != "jane@example.com" || created.identity.EmailVerified {
		t.Fatal("incorrect registered user or identity")
	}
	if created.identity.PasswordHash == nil || bcrypt.CompareHashAndPassword([]byte(*created.identity.PasswordHash), []byte(request.Password)) != nil {
		t.Fatal("password was not stored as a valid bcrypt hash")
	}
	if created.user.AvatarURL == nil || !strings.HasPrefix(*created.user.AvatarURL, storage.AvatarPrefix) || strings.Contains(*created.user.AvatarURL, "original") {
		t.Fatal("expected a generated relative avatar path")
	}
	file, err := os.Open(filepath.Join(directory, filepath.Base(*created.user.AvatarURL)))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := png.Decode(file); err != nil {
		t.Fatal("stored file is not a valid PNG")
	}
	if userID, err := ValidateJWT(response.Token); err != nil || userID != created.user.ID {
		t.Fatal("expected a token for the registered user")
	}
	for _, user := range []models.UserResponse{response.User, currentUser(t, service, created.user.ID)} {
		if user.Email == nil || *user.Email != "jane@example.com" || user.ProfilePicture == nil || *user.ProfilePicture != "https://assets.example.com"+*created.user.AvatarURL {
			t.Fatalf("unexpected user profile: %+v", user)
		}
	}
	body, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), request.Password) || strings.Contains(string(body), "password_hash") || strings.Contains(string(body), "confirm_password") {
		t.Fatal("registration response leaked credentials")
	}
}

func currentUser(t *testing.T, service AuthService, userID uint) models.UserResponse {
	t.Helper()
	response, err := service.GetCurrentUser(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	return response.User
}

func TestRegisterWithoutAvatar(t *testing.T) {
	repository := &fakeUserRepository{}
	response, err := NewAuthService(repository).Register(context.Background(), registerRequest())
	if err != nil {
		t.Fatal(err)
	}
	if response.User.ProfilePicture != nil || repository.createdUsersWithIdentity[0].user.AvatarURL != nil {
		t.Fatal("expected no avatar")
	}
}

func TestRegisterDeletesAvatarOnDatabaseFailure(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want error
	}{
		{"database failure", errors.New("database offline"), nil},
		{"concurrent duplicate", repositories.ErrDuplicateIdentity, ErrDuplicateEmail},
	} {
		t.Run(tc.name, func(t *testing.T) {
			directory := t.TempDir()
			repository := &fakeUserRepository{createError: tc.err}
			service := NewAuthService(repository).(*authService)
			service.avatarStore = &storage.LocalAvatarStore{Directory: directory}
			request := registerRequest()
			request.ProfilePicture = registerAvatar(t)
			response, err := service.Register(context.Background(), request)
			if response != nil || err == nil || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Fatalf("unexpected registration result: %v, %v", response, err)
			}
			files, err := os.ReadDir(directory)
			if err != nil || len(files) != 0 {
				t.Fatal("failed registration left an orphaned avatar")
			}
			if len(repository.createdUsersWithIdentity) != 0 {
				t.Fatal("failed registration created a user")
			}
		})
	}
}

func TestRegisterValidatesBeforeRepositoryAccess(t *testing.T) {
	for _, tc := range []struct {
		name   string
		modify func(*models.RegisterRequest)
	}{
		{"blank display name", func(r *models.RegisterRequest) { r.DisplayName = "   " }},
		{"invalid email", func(r *models.RegisterRequest) { r.Email = "invalid" }},
		{"short password", func(r *models.RegisterRequest) { r.Password, r.ConfirmPassword = "short", "short" }},
		{"missing confirmation", func(r *models.RegisterRequest) { r.ConfirmPassword = "" }},
		{"mismatch", func(r *models.RegisterRequest) { r.ConfirmPassword = "different" }},
		{"bcrypt byte limit", func(r *models.RegisterRequest) {
			r.Password = strings.Repeat("ก", 25)
			r.ConfirmPassword = r.Password
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := registerRequest()
			tc.modify(&request)
			_, err := NewAuthService(nil).Register(context.Background(), request)
			var validationError *RegistrationValidationError
			if !errors.As(err, &validationError) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}
