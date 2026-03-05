package services

import "context"

type MockStorage struct {
	GetAliasByURLFn func(ctx context.Context, originalURL string) (string, error)
	CreateFn        func(ctx context.Context, alias, originalURL string) error
	GetURLFn        func(ctx context.Context, alias string) (string, error)
}

func (m *MockStorage) GetAliasByURL(ctx context.Context, originalURL string) (string, error) {
	return m.GetAliasByURLFn(ctx, originalURL)
}

func (m *MockStorage) Create(ctx context.Context, alias, originalURL string) error {
	return m.CreateFn(ctx, alias, originalURL)
}
func (m *MockStorage) GetURL(ctx context.Context, alias string) (string, error) {
	return m.GetURLFn(ctx, alias)
}
