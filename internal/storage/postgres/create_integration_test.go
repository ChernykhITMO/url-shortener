//go:build integration
// +build integration

package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func TestCreate_NewURL_ReturnsAlias(t *testing.T) {
	cleanupURLTable(t)

	ctx := context.Background()
	alias, err := testStore.Create(ctx, "fixedAlias", "https://example.com")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if alias != "fixedAlias" {
		t.Fatalf("expected alias %q, got %q", "fixedAlias", alias)
	}
}

func TestCreate_SameURL_ReturnsExistingAlias(t *testing.T) {
	cleanupURLTable(t)

	ctx := context.Background()
	firstAlias, err := testStore.Create(ctx, "FirstAlias", "https://same-url.com")
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	secondAlias, err := testStore.Create(ctx, "OtherAlias", "https://same-url.com")
	if err != nil {
		t.Fatalf("second Create failed: %v", err)
	}

	if secondAlias != firstAlias {
		t.Fatalf("expected alias %q, got %q", firstAlias, secondAlias)
	}
}

func TestCreate_AliasTaken_ReturnsErrAliasConflict(t *testing.T) {
	cleanupURLTable(t)

	ctx := context.Background()
	if _, err := testStore.Create(ctx, "fixedAlias", "https://first.com"); err != nil {
		t.Fatalf("seed Create failed: %v", err)
	}

	if _, err := testStore.Create(ctx, "fixedAlias", "https://second.com"); !errors.Is(err, storage.ErrAliasConflict) {
		t.Fatalf("expected ErrAliasConflict, got %v", err)
	}
}

func TestCreate_InvalidAlias_ReturnsErrInvalidAlias(t *testing.T) {
	cleanupURLTable(t)

	_, err := testStore.Create(context.Background(), "invalid", "https://example.com")
	if !errors.Is(err, storage.ErrInvalidAlias) {
		t.Fatalf("expected ErrInvalidAlias, got %v", err)
	}
}
