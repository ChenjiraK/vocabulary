package route

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLevelRouteRequiresAuthentication(t *testing.T) {
	router := gin.New()
	RegisterAPIRoutes(router, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/level", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want 401; body: %s", recorder.Code, recorder.Body.String())
	}
}

func TestAvatarRouteUsesConfiguredDirectory(t *testing.T) {
	directory := t.TempDir()
	t.Setenv("AVATAR_UPLOAD_DIR", directory)
	if err := os.WriteFile(filepath.Join(directory, "avatar.png"), []byte("avatar contents"), 0600); err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	RegisterAPIRoutes(router, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/uploads/avatars/avatar.png", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != "avatar contents" {
		t.Fatalf("unexpected avatar response: %d %s", recorder.Code, recorder.Body.String())
	}
}
