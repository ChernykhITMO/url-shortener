package router

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/http/handlers"
	"github.com/ChernykhITMO/url-shortener/internal/services"
	"github.com/ChernykhITMO/url-shortener/internal/storage/inmemory"
)

func TestRouter_CreateAliasThenGetURL(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	service := services.New(inmemory.New(), 20, 10)
	handler := handlers.New(log, service, 1024)

	server := httptest.NewServer(New(handler, log))
	defer server.Close()

	createReqBody := []byte(`{"url":"https://example.com"}`)
	createResp, err := http.Post(server.URL+"/url", "application/json", bytes.NewReader(createReqBody))
	if err != nil {
		t.Fatalf("create alias request failed: %v", err)
	}
	defer func() {
		_ = createResp.Body.Close()
	}()

	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createResp.StatusCode)
	}

	var createPayload struct {
		Alias string `json:"alias"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&createPayload); err != nil {
		t.Fatalf("decode create response failed: %v", err)
	}
	if createPayload.Alias == "" {
		t.Fatal("expected alias in create response")
	}

	getResp, err := http.Get(server.URL + "/url/" + createPayload.Alias)
	if err != nil {
		t.Fatalf("get url request failed: %v", err)
	}
	defer func() {
		_ = getResp.Body.Close()
	}()

	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, getResp.StatusCode)
	}

	var getPayload struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(getResp.Body).Decode(&getPayload); err != nil {
		t.Fatalf("decode get response failed: %v", err)
	}
	if getPayload.URL != "https://example.com/" {
		t.Fatalf("expected normalized url %q, got %q", "https://example.com/", getPayload.URL)
	}
}
