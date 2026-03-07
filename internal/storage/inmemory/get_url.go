package inmemory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetURL(_ context.Context, alias string) (string, error) {
	const op = "inmemory.GetURL"

	s.mux.RLock()
	url, ok := s.aliasToURL[alias]
	s.mux.RUnlock()

	if ok {
		return url, nil
	}

	return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
}
