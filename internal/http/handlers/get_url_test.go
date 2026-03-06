package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/services"
)

func TestGetURL_Success_Returns200AndURL(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	h := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			if alias != "fixedAlias" {
				t.Fatalf("unexpected alias: %s", alias)
			}
			return "https://google.com", nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/url/fixedAlias", nil)
	req.SetPathValue("alias", "fixedAlias")
	rec := httptest.NewRecorder()

	h.GetURL(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.URL != "https://google.com" {
		t.Fatalf("expected url %q, got %q", "https://google.com", resp.URL)
	}
}

func TestGetURL_InvalidAlias_Returns400(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	h := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			t.Fatal("service should not be called for invalid alias")
			return "", nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/url/short", nil)
	req.SetPathValue("alias", "short")
	rec := httptest.NewRecorder()

	h.GetURL(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.Error != msgInvalidAlias {
		t.Fatalf("expected error %q, got %q", msgInvalidAlias, resp.Error)
	}
}

func TestGetURL_NotFound_Returns404(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	h := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", services.ErrNotFound
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/url/fixedAlias", nil)
	req.SetPathValue("alias", "fixedAlias")
	rec := httptest.NewRecorder()

	h.GetURL(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.Error != msgNotFound {
		t.Fatalf("expected error %q, got %q", msgNotFound, resp.Error)
	}
}

func TestGetURL_ServiceUnknownError_Returns500(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	h := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", errors.New("db down")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/url/fixedAlias", nil)
	req.SetPathValue("alias", "fixedAlias")
	rec := httptest.NewRecorder()

	h.GetURL(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.Error != msgInternalError {
		t.Fatalf("expected error %q, got %q", msgInternalError, resp.Error)
	}
}
