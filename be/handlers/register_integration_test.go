package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"vocabulary/middlewares"
	"vocabulary/models"
	"vocabulary/repositories"
	"vocabulary/services"
	"vocabulary/storage"

	"github.com/gin-gonic/gin"
)

// This test exercises HTTP binding, the real service, image storage, JWT auth,
// and profile rendering; only database persistence is replaced.
func TestRegisterThenFetchCurrentUserAndAvatar(t *testing.T) {
	t.Setenv("AVATAR_UPLOAD_DIR", t.TempDir())
	t.Setenv("PUBLIC_ASSET_BASE_URL", "https://api.example.com")
	repository := &registerRepositoryStub{}
	handler := NewAuthHandler(services.NewAuthService(repository))
	router := gin.New()
	router.POST("/auth/register", handler.Register)
	router.GET("/auth/me", middlewares.AuthMiddleware(), handler.GetCurrentUser)
	router.StaticFS("/uploads/avatars", gin.Dir(storage.AvatarDirectory(), false))

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range map[string]string{
		"display_name": "Jane", "email": "JANE@example.com",
		"password": "password123", "confirm_password": "password123",
	} {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("profile_picture", "picture.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(part, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/auth/register", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", recorder.Code, recorder.Body.String())
	}
	var registered models.AuthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &registered); err != nil {
		t.Fatal(err)
	}

	request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	request.Header.Set("Authorization", "Bearer "+registered.Token)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("current user: %d %s", recorder.Code, recorder.Body.String())
	}
	var current models.AuthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &current); err != nil {
		t.Fatal(err)
	}
	if current.User.DisplayName != "Jane" || current.User.Email == nil || *current.User.Email != "jane@example.com" || current.User.ProfilePicture == nil || registered.User.ProfilePicture == nil || *current.User.ProfilePicture != *registered.User.ProfilePicture {
		t.Fatalf("unexpected current user: %+v", current.User)
	}
	if repository.user.AvatarURL == nil {
		t.Fatal("avatar path was not persisted")
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, *repository.user.AvatarURL, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("avatar: %d", recorder.Code)
	}
	if _, err := png.Decode(bytes.NewReader(recorder.Body.Bytes())); err != nil {
		t.Fatal("avatar URL did not serve an image")
	}
}

type registerRepositoryStub struct {
	repositories.UserRepository
	user     models.User
	identity models.UserAuthIdentity
}

func (r *registerRepositoryStub) FindIdentityByEmail(context.Context, string) (*models.UserAuthIdentity, error) {
	return nil, nil
}

func (r *registerRepositoryStub) CreateUserWithIdentity(_ context.Context, user *models.User, identity *models.UserAuthIdentity) error {
	user.ID = 7
	identity.UserID = user.ID
	r.user, r.identity = *user, *identity
	return nil
}

func (r *registerRepositoryStub) FindUserByID(_ context.Context, id uint) (*models.User, error) {
	if id != r.user.ID {
		return nil, nil
	}
	return &r.user, nil
}

func (r *registerRepositoryStub) FindEmailByUserID(context.Context, uint) (*string, error) {
	return r.identity.Email, nil
}

func (r *registerRepositoryStub) GetUserGameStats(context.Context, uint) (models.UserGameStats, error) {
	return models.UserGameStats{}, nil
}
