package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func TestGetURL_Success_ReturnsURL(t *testing.T) {
	url := "https://google.com"
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		CreateFn:        func(ctx context.Context, alias, originalURL string) error { return nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "https://google.com", nil
		},
	}, 20)

	alias := "fixedAlias"
	got, err := s.GetURL(context.Background(), alias)
	if err != nil {
		t.Fatalf("expected %s, got err %v", url, err)
	}

	if got != url {
		t.Fatalf("expected %s, got %s", url, got)
	}
}

func TestGetURL_NotFound_ReturnsErrNotFound(t *testing.T) {
	alias := "fixedAlias"
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		CreateFn:        func(ctx context.Context, alias, originalURL string) error { return nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", storage.ErrNotFound
		},
	}, 20)

	_, err := s.GetURL(context.Background(), alias)
	if err == nil {
		t.Fatalf("expected %v, got nil", ErrNotFound)
	}

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected %v, got %v", ErrNotFound, err)
	}
}

func TestGetURL_StorageError_ReturnsError(t *testing.T) {
	errDBDown := errors.New("db down")
	alias := "fixedAlias"
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) { return "", nil },
		CreateFn:        func(ctx context.Context, alias, originalURL string) error { return nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", errDBDown
		},
	}, 20)

	_, err := s.GetURL(context.Background(), alias)
	if err == nil {
		t.Fatalf("expected %v, got nil", errDBDown)
	}

	if errors.Is(err, ErrNotFound) {
		t.Fatalf("expected %v, got not found", errDBDown)
	}
	if !errors.Is(err, errDBDown) {
		t.Fatalf("expected %v, got %v", errDBDown, err)
	}
}
