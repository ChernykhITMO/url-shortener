package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
)

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "postgres.GetURL"

	const query = `SELECT original_url FROM url WHERE alias = $1`

	var originalURL string
	if err := s.db.QueryRowContext(ctx, query, alias).Scan(&originalURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return originalURL, nil
}
