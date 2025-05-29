package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// setupVersionTest creates a new test environment with SystemAPIController
func setupVersionTest(t *testing.T) (*SystemAPIService, *SystemAPIController, *mux.Router, *config.Config, time.Time) {
	buildDate := time.Now().UTC()
	cfg := &config.Config{Version: "vX", BuildDate: buildDate, CommitHash: "hX"}
	svc := NewSystemAPIService(cfg)
	ctrl := NewSystemAPIController(svc)
	router := NewRouter(ctrl)
	if router == nil {
		t.Fatal("router is nil")
	}
	return svc, ctrl, router, cfg, buildDate
}

func TestGetVersionService(t *testing.T) {
	svc, _, _, cfg, buildDate := setupVersionTest(t)

	res, err := svc.GetVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Code != http.StatusOK {
		t.Fatalf("expected code 200, got %d", res.Code)
	}
	v, ok := res.Body.(models.Version)
	if !ok {
		t.Fatalf("expected models.Version, got %T", res.Body)
	}
	if v.Version != cfg.Version {
		t.Errorf("expected version %s, got %s", cfg.Version, v.Version)
	}
	if !v.BuildDate.Equal(buildDate) {
		t.Errorf("expected build date %v, got %v", buildDate, v.BuildDate)
	}
	if v.CommitHash != cfg.CommitHash {
		t.Errorf("expected commit hash %s, got %s", cfg.CommitHash, v.CommitHash)
	}
}

func TestGetVersionEndpoint(t *testing.T) {
	_, ctrl, _, cfg, buildDate := setupVersionTest(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	ctrl.GetVersion(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status code %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("invalid content-type, got %s", ct)
	}
	var data models.Version
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if data.Version != cfg.Version {
		t.Errorf("expected %s, got %s", cfg.Version, data.Version)
	}
	if !data.BuildDate.Equal(buildDate) {
		t.Errorf("expected build date %v, got %v", buildDate, data.BuildDate)
	}
	if data.CommitHash != cfg.CommitHash {
		t.Errorf("expected %s, got %s", cfg.CommitHash, data.CommitHash)
	}
}

func TestGetVersionEndpointInvalidMethods(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	_, _, router, _, _ := setupVersionTest(t)

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/version", nil)
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status code %d for method %s, got %d",
					http.StatusMethodNotAllowed, method, rr.Code)
			}
		})
	}
}

func TestGetVersionEndpointIgnoreHeaders(t *testing.T) {
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

	_, _, router, _, buildDate := setupVersionTest(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/version", nil)
			tc.setupRequest(req)

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			contentType := rr.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				t.Errorf("invalid content-type, got %s", contentType)
			}

			var data models.Version
			if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}
			if data.Version != "vX" {
				t.Errorf("expected version vX, got %s", data.Version)
			}
			if !data.BuildDate.Equal(buildDate) {
				t.Errorf("expected build date to be close to now, got %v", data.BuildDate)
			}
			if data.CommitHash != "hX" {
				t.Errorf("expected commit hash hX, got %s", data.CommitHash)
			}
		})
	}
}
