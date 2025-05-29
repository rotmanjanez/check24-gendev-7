package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// setupHealthTest creates a new test environment with HealthAPIController and router
func setupHealthTest(t *testing.T) (*HealthAPIService, *HealthAPIController, *mux.Router) {
	svc := NewHealthAPIService()
	ctrl := NewHealthAPIController(svc)
	router := NewRouter(ctrl)
	if router == nil {
		t.Fatal("router is nil")
	}
	return svc, ctrl, router
}

func TestHealthCheckService(t *testing.T) {
	svc, _, _ := setupHealthTest(t)
	res, err := svc.HealthCheck(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", res.Code)
	}
	h, ok := res.Body.(models.Health)
	if !ok {
		t.Fatalf("expected models.Health, got %T", res.Body)
	}
	if h.Status != "ok" {
		t.Errorf("expected status ok, got %s", h.Status)
	}
}

func TestHealthCheckEndpoint(t *testing.T) {
	_, ctrl, _ := setupHealthTest(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	ctrl.HealthCheck(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status code %d", rr.Code)
	}
	var data models.Health
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data.Status != "ok" {
		t.Errorf("expected status ok, got %s", data.Status)
	}
}

func TestHealthCheckEndpointInvalidMethods(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	_, _, router := setupHealthTest(t)
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/health", nil)
			router.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status code %d for method %s, got %d",
					http.StatusMethodNotAllowed, method, rr.Code)
			}
		})
	}
}

func TestHealthCheckEndpointIgnoreHeaders(t *testing.T) {
	testCases := []struct {
		name         string
		setupRequest func(*http.Request)
	}{
		{
			name: "invalid content type",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/xml")
			},
		},
		{
			name: "invalid accept header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Accept", "application/xml")
			},
		},
		{
			name: "with request body",
			setupRequest: func(req *http.Request) {
				req.Body = io.NopCloser(strings.NewReader(`{"invalid": "payload"}`))
				req.Header.Set("Content-Type", "application/json")
			},
		},
	}

	_, _, router := setupHealthTest(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			tc.setupRequest(req)
			router.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}
			contentType := rr.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				t.Errorf("invalid content-type, got %s", contentType)
			}
			var data models.Health
			if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}
			if data.Status != "ok" {
				t.Errorf("expected status ok, got %s", data.Status)
			}
		})
	}
}
