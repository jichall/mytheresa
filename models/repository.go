package models

import "context"

type RepositoryObject interface{ Product | Category }

type Repository[T RepositoryObject] interface {
	GetAll(ctx context.Context) ([]T, error)
	GetPaged(ctx context.Context, page int, size int) ([]T, error)
	GetByCode(ctx context.Context, code string) (*T, error)

	Save(ctx context.Context, data *T) error
}
