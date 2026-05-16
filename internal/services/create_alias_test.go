package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func TestCreateAlias_InvalidURL_ReturnsErrInvalidURL(t *testing.T) {
	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) { return "", nil },
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
		"ftp://example.com",
	}

	for _, data := range testData {
		t.Run(data, func(t *testing.T) {
			_, err := s.CreateAlias(context.Background(), data, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !(errors.Is(err, ErrInvalidURL)) {
				t.Fatalf("expected ErrInvalidURL, got %v", err)
			}
		})
	}
}

func TestCreateAlias_NormalizesURLBeforeStore(t *testing.T) {
	var gotURL string

	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			gotURL = originalURL
			return alias, nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	}, 20)

	_, err := s.CreateAlias(context.Background(), "  HTTPS://Example.COM:443  ", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotURL != "https://example.com/" {
		t.Fatalf("expected normalized url %q, got %q", "https://example.com/", gotURL)
	}
}

func TestCreateAlias_ValidURL_CreatesAlias(t *testing.T) {
	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) { return alias, nil },
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
			_, err := s.CreateAlias(context.Background(), data, "")
			if err != nil {
				t.Fatalf("expected correct function, got error %v", err)
			}
		})
	}
}

func TestCreateAlias_AliasConflict_ReturnsRetrySuccess(t *testing.T) {
	var createCalls, aliasCalls int

	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			createCalls++
			if createCalls == 1 {
				return "", storage.ErrAliasConflict
			}
			return alias, nil
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

	got, err := s.CreateAlias(context.Background(), url, "")
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
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			createCalls++
			return "fixedAlias", nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"
	expectedAlias := "fixedAlias"

	alias, err := s.CreateAlias(context.Background(), url, "")
	if err != nil {
		t.Fatalf("expected nil, got err %v", err)
	}

	if alias != expectedAlias {
		t.Fatalf("expected %s, got %s", expectedAlias, alias)
	}

	if createCalls != 1 {
		t.Fatalf("expected Create to be called 1 time, got %d", createCalls)
	}
}

func TestCreateAlias_StorageError_ReturnsError(t *testing.T) {
	errDBDown := errors.New("db down")
	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			return "", errDBDown
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"

	_, err := s.CreateAlias(context.Background(), url, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errDBDown) {
		t.Fatalf("expected %v, got %v", errDBDown, err)
	}
}

func TestCreateAlias_AttemptsExceeded_ReturnsErrAttemptsExceeded(t *testing.T) {
	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			return "", storage.ErrAliasConflict
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) {
			return "", nil
		},
	}, 20)

	url := "https://google.com"

	alias, err := s.CreateAlias(context.Background(), url, "")
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

func TestCreateAlias_CustomAlias_UsesProvidedAlias(t *testing.T) {
	var gotAlias string

	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			gotAlias = alias
			return alias, nil
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	}, 20)

	alias, err := s.CreateAlias(context.Background(), "https://google.com", "MyAlias_01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alias != "MyAlias_01" || gotAlias != "MyAlias_01" {
		t.Fatalf("expected custom alias to be used, got result=%q stored=%q", alias, gotAlias)
	}
}

func TestCreateAlias_CustomAliasConflict_ReturnsErrAliasTaken(t *testing.T) {
	s := New(&MockStorage{
		CreateFn: func(ctx context.Context, alias, originalURL string) (string, error) {
			return "", storage.ErrAliasConflict
		},
		GetURLFn: func(ctx context.Context, alias string) (string, error) { return "", nil },
	}, 20)

	_, err := s.CreateAlias(context.Background(), "https://google.com", "MyAlias_01")
	if !errors.Is(err, ErrAliasTaken) {
		t.Fatalf("expected %v, got %v", ErrAliasTaken, err)
	}
}
