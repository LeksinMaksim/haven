package todo

import "context"

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Todo, error)
	GetByID(ctx context.Context, id int64) (*Todo, error)
	List(ctx context.Context, filter ListFilter) ([]Todo, error)
	Update(ctx context.Context, id int64, req UpdateRequest) (*Todo, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*Todo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	t := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*Todo, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, filter ListFilter) ([]Todo, error) {
	todos, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	if todos == nil {
		todos = []Todo{}
	}
	return todos, nil
}

func (s *service) Update(ctx context.Context, id int64, req UpdateRequest) (*Todo, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		t.Title = *req.Title
	}
	if req.Description != nil {
		t.Description = *req.Description
	}
	if req.Done != nil {
		t.Done = *req.Done
	}
	if req.Priority != nil {
		t.Priority = *req.Priority
	}
	if req.DueDate != nil {
		t.DueDate = req.DueDate
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, err
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
