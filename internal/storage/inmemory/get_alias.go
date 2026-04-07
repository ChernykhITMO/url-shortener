package inmemory

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetAlias(_ context.Context, originalURL string) (string, error) {
	const op = "inmemory.GetAlias"

	s.mux.RLock()
	alias, ok := s.urlToAlias[originalURL]
	s.mux.RUnlock()

	if ok {
		return alias, nil
	}

	return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
}
