package middleware

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireMethodMiddleware_WrongMethod_Returns405(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not be called")
	})

	h := RequireMethodMiddleware(http.MethodPost, log)(next)

	req := httptest.NewRequest(http.MethodGet, "/url", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("expected Allow=%q, got %q", http.MethodPost, got)
	}

	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Error != "method not allowed" {
		t.Fatalf("unexpected error: %q", resp.Error)
	}

}
