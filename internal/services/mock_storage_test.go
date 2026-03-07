package services

import "context"

type MockStorage struct {
	CreateFn func(ctx context.Context, alias, originalURL string) (string, error)
	GetURLFn func(ctx context.Context, alias string) (string, error)
}

func (m *MockStorage) Create(ctx context.Context, alias, originalURL string) (string, error) {
	return m.CreateFn(ctx, alias, originalURL)
}
func (m *MockStorage) GetURL(ctx context.Context, alias string) (string, error) {
	return m.GetURLFn(ctx, alias)
}
