package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/url"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	lengthAlias = 10
)

type Storage interface {
	GetAliasByURL(ctx context.Context, originalURL string) (string, error)
	Create(ctx context.Context, alias, originalURL string) error
	GetURL(ctx context.Context, alias string) (string, error)
}

type Service struct {
	storage       Storage
	maxAttempts   int
	generateAlias func() (string, error)
}

func New(storage Storage, maxAttempts int) *Service {
	return &Service{
		storage:       storage,
		maxAttempts:   maxAttempts,
		generateAlias: generateAlias,
	}
}

func validateURL(row string) error {
	const op = "services.validateURL"
	u, err := url.ParseRequestURI(row)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return wrapError(op, ErrInvalidURL)
	}
	return nil
}

func generateAlias() (string, error) {
	const op = "services.genAlias"
	b := make([]byte, lengthAlias)

	if _, err := rand.Read(b); err != nil {
		return "", wrapError(op, err)
	}

	for i := 0; i < lengthAlias; i++ {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}

	return string(b), nil
}

func wrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
