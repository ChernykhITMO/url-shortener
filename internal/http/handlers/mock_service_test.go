package handlers

import "context"

type MockService struct {
	CreateAliasFn func(ctx context.Context, originalURL string) (string, error)
	GetURLFn      func(ctx context.Context, alias string) (string, error)
}

func (m *MockService) CreateAlias(ctx context.Context, originalURL string) (string, error) {
	return m.CreateAliasFn(ctx, originalURL)
}

func (m *MockService) GetURL(ctx context.Context, alias string) (string, error) {
	return m.GetURLFn(ctx, alias)
}
