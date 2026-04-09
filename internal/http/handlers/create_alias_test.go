package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type errorResponse struct {
	Error string `json:"error"`
}

func TestCreateAlias_Success_Returns201AndAlias(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			if originalURL != "https://google.com" {
				t.Fatalf("unexpected url: %s", originalURL)
			}
			if requestedAlias != "" {
				t.Fatalf("unexpected alias: %s", requestedAlias)
			}
			return "fixedAlias", nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	})

	body := []byte(`{"url": "https://google.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp struct {
		Alias string `json:"alias"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp.Alias != "fixedAlias" {
		t.Fatalf("expected alias %q, got %q", "fixedAlias", resp.Alias)
	}
	if location := rec.Header().Get("Location"); location != "/url/fixedAlias" {
		t.Fatalf("expected location %q, got %q", "/url/fixedAlias", location)
	}
}

func TestCreateAlias_InvalidJSON_Returns400AndErrorJSON(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			t.Fatal("service should not be called for invalid json")
			return "", nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	})

	body := []byte(`{"url": "https://google.com"`)
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.Error != msgInvalidJSON {
		t.Fatalf("expected error %q, got %q", msgInvalidJSON, resp.Error)
	}
}

func TestCreateAlias_PayloadTooLarge_Returns413(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			return "", nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	})

	tooLargeURL := strings.Repeat("a", maxCreateAliasBodyBytes+100)
	body := []byte(`{"url":"` + tooLargeURL + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp.Error != msgPayloadTooLarge {
		t.Fatalf("expected error %q, got %q", msgPayloadTooLarge, resp.Error)
	}
}

func TestCreateAlias_ServiceUnknownError_Returns500(t *testing.T) {
	errDBDown := errors.New("db down")
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			return "", errDBDown
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	})

	body := []byte(`{"url":"https://google.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

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

func TestCreateAlias_CustomAlias_Returns201(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			if requestedAlias != "MyAlias_01" {
				t.Fatalf("unexpected alias: %s", requestedAlias)
			}
			return requestedAlias, nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(`{"url":"https://google.com","alias":"MyAlias_01"}`))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestCreateAlias_UnsupportedContentType_Returns415(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	s := New(log, &MockService{
		CreateAliasFn: func(ctx context.Context, originalURL, requestedAlias string) (string, error) {
			t.Fatal("service should not be called for unsupported content type")
			return "", nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewBufferString(`{"url":"https://google.com"}`))
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	s.CreateAlias(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status %d, got %d", http.StatusUnsupportedMediaType, rec.Code)
	}
}
