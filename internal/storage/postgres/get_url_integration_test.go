//go:build integration
// +build integration

package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func TestGetURL_NotFound_ReturnsErrNotFound(t *testing.T) {
	cleanupURLTable(t)

	ctx := context.Background()
	if _, err := testStore.GetURL(ctx, "UnknownAlia"); !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

