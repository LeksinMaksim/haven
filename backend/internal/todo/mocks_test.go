package todo

import (
	"context"
)

type mockRepository struct {
	createFunc  func(ctx context.Context, t *Todo) error
	getByIDFunc func(ctx context.Context, id int64) (*Todo, error)
	listFunc    func(ctx context.Context, filter ListFilter) ([]Todo, error)
	updateFunc  func(ctx context.Context, t *Todo) error
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockRepository) Create(ctx context.Context, t *Todo) error {
	return m.createFunc(ctx, t)
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*Todo, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockRepository) List(ctx context.Context, filter ListFilter) ([]Todo, error) {
	return m.listFunc(ctx, filter)
}

func (m *mockRepository) Update(ctx context.Context, t *Todo) error {
	return m.updateFunc(ctx, t)
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

type mockService struct {
	createFunc  func(ctx context.Context, req CreateRequest) (*Todo, error)
	getByIDFunc func(ctx context.Context, id int64) (*Todo, error)
	listFunc    func(ctx context.Context, filter ListFilter) ([]Todo, error)
	updateFunc  func(ctx context.Context, id int64, req UpdateRequest) (*Todo, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockService) Create(ctx context.Context, req CreateRequest) (*Todo, error) {
	return m.createFunc(ctx, req)
}

func (m *mockService) GetByID(ctx context.Context, id int64) (*Todo, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *mockService) List(ctx context.Context, filter ListFilter) ([]Todo, error) {
	return m.listFunc(ctx, filter)
}

func (m *mockService) Update(ctx context.Context, id int64, req UpdateRequest) (*Todo, error) {
	return m.updateFunc(ctx, id, req)
}

func (m *mockService) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}
