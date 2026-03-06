package in_memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestInMemory_CreateAndGet(t *testing.T) {
	s := New()

	url := "https://google.com"
	alias := "fixedAlias"

	ctx := context.Background()

	if err := s.Create(ctx, alias, url); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	gotURL, err := s.GetURL(ctx, alias)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if gotURL != url {
		t.Fatalf("expected %s, got %s", url, gotURL)
	}

	gotAlias, err := s.GetAliasByURL(ctx, url)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if gotAlias != alias {
		t.Fatalf("expected %s, got %s", alias, gotAlias)
	}
}

func TestInMemory_ConcurrentAccess(t *testing.T) {
	s := New()
	ctx := context.Background()

	const n = 1000

	var wg sync.WaitGroup
	wg.Add(n)
	errCh := make(chan error, n*3)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()

			alias := fmt.Sprintf("alias_%04d", i)
			url := fmt.Sprintf("https://example.com/%d", i)

			if err := s.Create(ctx, alias, url); err != nil {
				errCh <- fmt.Errorf("create failed for alias=%s url=%s: %w", alias, url, err)
			}
			if _, err := s.GetURL(ctx, alias); err != nil {
				errCh <- fmt.Errorf("get url failed for alias=%s: %w", alias, err)
			}
			if _, err := s.GetAliasByURL(ctx, url); err != nil {
				errCh <- fmt.Errorf("get alias failed for url=%s: %w", url, err)
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("unexpected error in concurrent access: %v", err)
	}
}
