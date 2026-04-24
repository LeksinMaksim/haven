package todo

import "context"

type Repository interface {
	Create(ctx context.Context, t *Todo) error
	GetByID(ctx context.Context, id int64) (*Todo, error)
	List(ctx context.Context, filter ListFilter) ([]Todo, error)
	Update(ctx context.Context, t *Todo) error
	Delete (ctx context.Context, id int64) error
}

type ListFilter struct {
	Done *bool
	Priority *Priority
}