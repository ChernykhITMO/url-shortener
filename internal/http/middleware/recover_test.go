package middleware

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware_Panic_Returns500(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	})

	h := RecoverMiddleware(log)(next)

	req := httptest.NewRequest(http.MethodGet, "/url/x", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Error != "internal error" {
		t.Fatalf("unexpected error: %q", resp.Error)
	}
}
