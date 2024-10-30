package router

import (
	"github.com/Ayano2000/push/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_TestHandleStaticRoute(t *testing.T) {
	router := NewDynamicMux(&handlers.Handler{})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("POST /webhooks", handler)

	req, _ := http.NewRequest("POST", "/webhooks", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRouter_HandleDynamicRouteParamsAreSet(t *testing.T) {
	// single param
	router := NewDynamicMux(&handlers.Handler{})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(urlParamContextKey).(map[string]string)
		if params["name"] != "asdf" {
			t.Errorf("expected parameter 'name' to be 'asdf', got %s", params["name"])
		}
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("GET /webhooks/{name}/content", handler)

	req, _ := http.NewRequest("GET", "/webhooks/asdf/content", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Multiple params
	router = NewDynamicMux(&handlers.Handler{})
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(urlParamContextKey).(map[string]string)
		if params["name"] != "asdf" {
			t.Errorf("expected parameter 'name' to be 'asdf', got %s", params["name"])
		}
		if params["id"] != "1234" {
			t.Errorf("expected parameter 'id' to be '1234', got %s", params["id"])
		}
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("PUT /webhooks/{name}/content/{id}", handler)

	req, _ = http.NewRequest("PUT", "/webhooks/asdf/content/1234", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

}

func TestRouter_MatchRoute(t *testing.T) {
	// Arrange
	route := &Route{
		pattern:    "PUT /webhooks/{name}/content/{id}",
		parameters: []string{"name", "id"},
		segments:   []string{"", "webhooks", "*", "content", "*"},
		isDynamic:  true,
	}

	// Act
	params, match := matchRoute(route, "PUT", "/webhooks/asdf/content/1234")

	// Assert
	if !match {
		t.Errorf("expected route to match, but it did not")
	}
	if params["name"] != "asdf" {
		t.Errorf("expected parameter 'name' to be 'asdf', got %s", params["name"])
	}
}
