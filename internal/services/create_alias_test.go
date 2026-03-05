package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func TestCreateAlias_InvalidURL_ReturnsErrInvalidURL(t *testing.T) {
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "", nil
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error { return nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	}, 20)

	testData := []string{
		"google.com",
		"http://",
		"https://",
		"",
		"    ",
		"my.itmo.ru",
		"https://goog le.com",
		"not_a_url",
		"not a url",
		"https://?q=1",
	}

	for _, data := range testData {
		t.Run(data, func(t *testing.T) {
			_, err := s.CreateAlias(context.Background(), data)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !(errors.Is(err, ErrInvalidURL)) {
				t.Fatalf("expected ErrInvalidURL, got %v", err)
			}
		})
	}
}

func TestCreateAlias_ValidURL_CreatesAlias(t *testing.T) {
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "", storage.ErrNotFound
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error { return nil },
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	}, 20)

	testData := []string{
		"http://google.com",
		"http://vk.com/feed",
		"https://my.itmo.ru",
		"https://google.com?x=1",
		"https://ozon.ru?x=1#1",
	}

	for _, data := range testData {
		t.Run(data, func(t *testing.T) {
			_, err := s.CreateAlias(context.Background(), data)
			if err != nil {
				t.Fatalf("expected correct function, got error %v", err)
			}
		})
	}
}

func TestCreateAlias_AliasConflict_ReturnsRetrySuccess(t *testing.T) {
	var createCalls, aliasCalls int

	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "", storage.ErrNotFound
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error {
			createCalls++
			if createCalls == 1 {
				return storage.ErrAliasConflict
			}
			return nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"

	s.generateAlias = func() (string, error) {
		aliasCalls++
		if aliasCalls == 1 {
			return "alias1", nil
		}
		return "alias2", nil
	}

	got, err := s.CreateAlias(context.Background(), url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "alias2" {
		t.Fatalf("expected alias2, got %s", got)
	}
	if createCalls != 2 {
		t.Fatalf("expected 2 create calls, got %d", createCalls)
	}
}

func TestCreateAlias_ExistingURL_ReturnsExistingAlias(t *testing.T) {
	createCalls := 0
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "fixedAlias", nil
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error {
			createCalls++
			return nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"
	expectedAlias := "fixedAlias"

	alias, err := s.CreateAlias(context.Background(), url)
	if err != nil {
		t.Fatalf("expected nil, got err %v", err)
	}

	if alias != expectedAlias {
		t.Fatalf("expected %s, got %s", expectedAlias, alias)
	}

	if createCalls != 0 {
		t.Fatalf("expected Create to be called 0 times, got %d", createCalls)
	}
}

func TestCreateAlias_URLConflict_ReturnsExistingAlias(t *testing.T) {
	getAliasCalls := 0
	expectedAlias := "fixedAlias"
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			getAliasCalls++
			if getAliasCalls == 1 {
				return "", storage.ErrNotFound
			}
			return expectedAlias, nil
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error {
			return storage.ErrURLConflict
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"

	alias, err := s.CreateAlias(context.Background(), url)
	if err != nil {
		t.Fatalf("expected nil, got err %v", err)
	}

	if alias != expectedAlias {
		t.Fatalf("expected %s, got %s", expectedAlias, alias)
	}

	if getAliasCalls != 2 {
		t.Fatalf("expected GetAliasByURL be called 2 times, got %d", getAliasCalls)
	}
}

func TestCreateAlias_AttemptsExceeded_ReturnsErrAttemptsExceeded(t *testing.T) {
	s := New(&MockStorage{
		GetAliasByURLFn: func(ctx context.Context, originalURL string) (string, error) {
			return "", storage.ErrNotFound
		},
		CreateFn: func(ctx context.Context, alias, originalURL string) error {
			return storage.ErrAliasConflict
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"

	alias, err := s.CreateAlias(context.Background(), url)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrAttemptsExceeded) {
		t.Fatalf("expected %v, got %v", ErrAttemptsExceeded, err)
	}

	if alias != "" {
		t.Fatalf("expected empty alias, got %v", alias)
	}
}
