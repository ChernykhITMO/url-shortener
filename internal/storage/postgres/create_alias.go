package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ChernykhITMO/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolation = "23505"
	checkViolation  = "23514"

	constraintPkAlias          = "url_pkey"
	constraintAliasLengthCheck = "url_alias_length_check"
	constraintAliasFormatCheck = "url_alias_format_check"

	attempts = 3
	delay    = 10 * time.Millisecond
)

func (s *Storage) Create(ctx context.Context, alias, originalURL string) (string, error) {
	const op = "postgres.Create"

	const cteQuery = `
  WITH ins AS (
        INSERT INTO url(alias, original_url)
        VALUES ($1, $2)
        ON CONFLICT (original_url) DO NOTHING
        RETURNING alias
  )
  SELECT COALESCE(
        (SELECT alias FROM ins),
        (SELECT alias FROM url WHERE original_url = $2)
  );
  `

	var gotAlias string
	err := s.db.QueryRowContext(ctx, cteQuery, alias, originalURL).Scan(&gotAlias)
	if err == nil {
		return gotAlias, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return s.selectAliasWithRetry(ctx, originalURL, op)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == uniqueViolation && pgErr.ConstraintName == constraintPkAlias {
			return "", fmt.Errorf("%s: %w", op, storage.ErrAliasConflict)
		}
		if pgErr.Code == checkViolation &&
			(pgErr.ConstraintName == constraintAliasLengthCheck || pgErr.ConstraintName == constraintAliasFormatCheck) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrInvalidAlias)
		}
	}

	return "", fmt.Errorf("%s: %w", op, err)
}

func (s *Storage) selectAliasWithRetry(ctx context.Context, originalURL, op string) (string, error) {
	const selectQuery = `SELECT alias FROM url WHERE original_url = $1`

	var alias string
	for i := 0; i < attempts; i++ {
		err := s.db.QueryRowContext(ctx, selectQuery, originalURL).Scan(&alias)
		if err == nil {
			return alias, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("%s: %w", op, ctx.Err())
		case <-time.After(delay):
		}
	}

	return "", fmt.Errorf("%s: fallback select: %w", op, storage.ErrNotFound)
}
