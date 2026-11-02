package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"vocabulary/models"
	"vocabulary/services"
	"vocabulary/storage"

	"github.com/gin-gonic/gin"
)

type registerServiceStub struct {
	services.AuthService
	called  bool
	request models.RegisterRequest
	err     error
}

func (s *registerServiceStub) Register(ctx context.Context, request models.RegisterRequest) (*models.AuthResponse, error) {
	s.called = true
	s.request = request
	return &models.AuthResponse{Message: "registered successfully"}, s.err
}

func TestRegisterConfirmPassword(t *testing.T) {
	for _, encoding := range []string{"json", "form", "multipart"} {
		for _, tc := range []struct {
			name    string
			confirm *string
			status  int
			message string
		}{
			{name: "matching", confirm: testStringPointer("password123"), status: http.StatusCreated, message: "registered successfully"},
			{name: "missing", status: http.StatusBadRequest, message: "confirm_password is required"},
			{name: "empty", confirm: testStringPointer(""), status: http.StatusBadRequest, message: "confirm_password is required"},
			{name: "different", confirm: testStringPointer("different123"), status: http.StatusBadRequest, message: "confirm_password must match password"},
			{name: "trailing space", confirm: testStringPointer("password123 "), status: http.StatusBadRequest, message: "confirm_password must match password"},
		} {
			t.Run(encoding+"/"+tc.name, func(t *testing.T) {
				fields := map[string]string{
					"display_name": "Jane",
					"email":        "jane@example.com",
					"password":     "password123",
				}
				if tc.confirm != nil {
					fields["confirm_password"] = *tc.confirm
				}

				var body bytes.Buffer
				var contentType string
				switch encoding {
				case "json":
					if err := json.NewEncoder(&body).Encode(fields); err != nil {
						t.Fatal(err)
					}
					contentType = "application/json"
				case "form":
					values := url.Values{}
					for key, value := range fields {
						values.Set(key, value)
					}
					body.WriteString(values.Encode())
					contentType = "application/x-www-form-urlencoded"
				case "multipart":
					writer := multipart.NewWriter(&body)
					for key, value := range fields {
						if err := writer.WriteField(key, value); err != nil {
							t.Fatal(err)
						}
					}
					if err := writer.Close(); err != nil {
						t.Fatal(err)
					}
					contentType = writer.FormDataContentType()
				}

				service := &registerServiceStub{}
				router := gin.New()
				router.POST("/auth/register", NewAuthHandler(service).Register)
				request := httptest.NewRequest(http.MethodPost, "/auth/register", &body)
				request.Header.Set("Content-Type", contentType)
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, request)

				if recorder.Code != tc.status {
					t.Fatalf("got status %d, want %d; body: %s", recorder.Code, tc.status, recorder.Body.String())
				}
				var response struct {
					Message string `json:"message"`
				}
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Message != tc.message {
					t.Fatalf("got message %q, want %q", response.Message, tc.message)
				}
				if service.called != (tc.status == http.StatusCreated) {
					t.Fatalf("unexpected service call: %t", service.called)
				}
				if service.called && (service.request.DisplayName != "Jane" || service.request.Email != "jane@example.com" || service.request.Password != "password123" || service.request.ConfirmPassword != "password123") {
					t.Fatal("registration fields were not bound correctly")
				}
				if strings.Contains(recorder.Body.String(), "password123") || strings.Contains(recorder.Body.String(), "different123") {
					t.Fatal("response exposed a password")
				}
			})
		}
	}
}

func testStringPointer(value string) *string {
	return &value
}

func TestRegisterErrorResponses(t *testing.T) {
	for _, tc := range []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{"duplicate email", services.ErrDuplicateEmail, http.StatusConflict, "email already exists"},
		{"invalid image", storage.ErrInvalidAvatar, http.StatusBadRequest, storage.ErrInvalidAvatar.Error()},
		{"large image", storage.ErrAvatarTooLarge, http.StatusRequestEntityTooLarge, storage.ErrAvatarTooLarge.Error()},
		{"validation", &services.RegistrationValidationError{Message: "password must be at most 72 bytes"}, http.StatusBadRequest, "password must be at most 72 bytes"},
		{"database error", errors.New("private database detail"), http.StatusInternalServerError, "cannot register user"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/auth/register", NewAuthHandler(&registerServiceStub{err: tc.err}).Register)
			request := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(`{"display_name":"Jane","email":"jane@example.com","password":"password123","confirm_password":"password123"}`))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			var body struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if recorder.Code != tc.status || body.Message != tc.message {
				t.Fatalf("unexpected response: %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestRegisterNormalizesFieldsAndRejectsMalformedRequests(t *testing.T) {
	for _, tc := range []struct {
		name, body, contentType string
		status                  int
	}{
		{"trimmed fields", `{"display_name":" Jane ","email":" JANE@Example.com ","password":"password123","confirm_password":"password123"}`, "application/json", http.StatusCreated},
		{"blank name", `{"display_name":"   ","email":"jane@example.com","password":"password123","confirm_password":"password123"}`, "application/json", http.StatusBadRequest},
		{"malformed JSON", `{`, "application/json", http.StatusBadRequest},
		{"unsupported content", `text`, "text/plain", http.StatusUnsupportedMediaType},
		{"oversized body", strings.Repeat("x", storage.MaxAvatarBytes+(64<<10)+1), "application/x-www-form-urlencoded", http.StatusRequestEntityTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &registerServiceStub{}
			router := gin.New()
			router.POST("/auth/register", NewAuthHandler(service).Register)
			request := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", tc.contentType)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != tc.status || service.called != (tc.status == http.StatusCreated) {
				t.Fatalf("unexpected response: %d %s", recorder.Code, recorder.Body.String())
			}
			if service.called && (service.request.DisplayName != "Jane" || service.request.Email != "jane@example.com") {
				t.Fatal("identity fields were not normalized")
			}
		})
	}
}
