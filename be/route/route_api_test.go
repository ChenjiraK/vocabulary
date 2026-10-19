package route

import (
	"net/http"
	"net/http/httptest"
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
