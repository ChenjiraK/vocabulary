package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"vocabulary/middlewares"
	"vocabulary/models"
	"vocabulary/services"

	"github.com/gin-gonic/gin"
)

type levelServiceStub struct {
	services.VocabularyService
	response models.VocabularyLevelsResponse
	err      error
	userID   uint
}

func (s *levelServiceStub) Levels(ctx context.Context, userID uint) (models.VocabularyLevelsResponse, error) {
	s.userID = userID
	return s.response, s.err
}

func TestGetLevels(t *testing.T) {
	want := models.VocabularyLevelsResponse{
		"A1": "inactive", "A2": "inactive", "B1": "complete",
		"B2": "active", "C1": "inactive", "C2": "inactive",
	}
	for _, tc := range []struct {
		name          string
		authenticated bool
		err           error
		status        int
	}{
		{name: "success", authenticated: true, status: http.StatusOK},
		{name: "unauthenticated", status: http.StatusUnauthorized},
		{name: "database failure", authenticated: true, err: errors.New("database failure"), status: http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &levelServiceStub{response: want, err: tc.err}
			router := gin.New()
			if tc.authenticated {
				router.Use(func(c *gin.Context) { c.Set(middlewares.UserIDContextKey, uint(42)) })
			}
			router.GET("/level", NewVocabularyHandler(service).GetLevels)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/level", nil))
			if recorder.Code != tc.status {
				t.Fatalf("got status %d, want %d", recorder.Code, tc.status)
			}
			if tc.authenticated && service.userID != 42 {
				t.Fatalf("expected user 42, got %d", service.userID)
			}
			if !tc.authenticated && service.userID != 0 {
				t.Fatal("service called without authentication")
			}
			if recorder.Code == http.StatusOK {
				var got models.VocabularyLevelsResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %v, want %v", got, want)
				}
			}
		})
	}
}
