package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/ChernykhITMO/url-shortener/internal/domain/alias"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type Storage interface {
	Create(ctx context.Context, alias, originalURL string) (string, error)
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

func normalizeAndValidateURL(rawURL string) (string, error) {
	const op = "services.normalizeAndValidateURL"

	trimmed := strings.TrimSpace(rawURL)
	u, err := url.ParseRequestURI(trimmed)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", wrapError(op, ErrInvalidURL)
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", wrapError(op, ErrInvalidURL)
	}

	hostname := strings.ToLower(u.Hostname())
	if hostname == "" {
		return "", wrapError(op, ErrInvalidURL)
	}

	port := u.Port()
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}

	host := hostname
	if port != "" {
		host = net.JoinHostPort(hostname, port)
	}

	u.Scheme = scheme
	u.Host = host
	if u.Path == "" {
		u.Path = "/"
	}

	return u.String(), nil
}

func generateAlias() (string, error) {
	const op = "services.genAlias"

	out := make([]byte, alias.Length)
	const byteRange = 256
	const maxNum = byteRange - (byteRange % len(alphabet))

	for i := 0; i < alias.Length; {
		var num [1]byte
		if _, err := rand.Read(num[:]); err != nil {
			return "", wrapError(op, err)
		}

		if int(num[0]) >= maxNum {
			continue
		}

		out[i] = alphabet[int(num[0])%len(alphabet)]
		i++
	}

	return string(out), nil
}

func wrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}

func validateAlias(s string) error {
	if len(s) != alias.Length {
		return ErrInvalidAlias
	}

	for i := 0; i < len(s); i++ {
		if s[i] >= 'a' && s[i] <= 'z' ||
			s[i] >= 'A' && s[i] <= 'Z' ||
			s[i] >= '0' && s[i] <= '9' ||
			s[i] == '_' {
			continue
		}
		return ErrInvalidAlias
	}
	return nil
}
