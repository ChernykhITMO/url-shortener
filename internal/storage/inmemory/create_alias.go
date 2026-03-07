package in_memory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) Create(ctx context.Context, alias, originalURL string) (string, error) {
	const op = "in_memory.Create"

	s.mux.Lock()
	defer s.mux.Unlock()

	if existingAlias, ok := s.urlToAlias[originalURL]; ok {
		return existingAlias, nil
	}

	if existingURL, ok := s.aliasToURL[alias]; ok && existingURL != originalURL {
		return "", fmt.Errorf("%s: %w", op, storage.ErrAliasConflict)
	}

	s.urlToAlias[originalURL] = alias
	s.aliasToURL[alias] = originalURL

	return alias, nil
}
