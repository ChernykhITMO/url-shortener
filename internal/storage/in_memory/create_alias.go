package in_memory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) Create(ctx context.Context, alias, originalURL string) error {
	const op = "in_memory.Create"

	s.mux.Lock()
	defer s.mux.Unlock()

	if existingURL, ok := s.aliasToURL[alias]; ok && existingURL != originalURL {
		return fmt.Errorf("%s: %w", op, storage.ErrAliasConflict)
	}
	if existingAlias, ok := s.urlToAlias[originalURL]; ok && existingAlias != alias {
		return fmt.Errorf("%s: %w", op, storage.ErrURLConflict)
	}

	s.urlToAlias[originalURL] = alias
	s.aliasToURL[alias] = originalURL

	return nil
}
