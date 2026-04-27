package todo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockRepository{
			createFunc: func(ctx context.Context, t *Todo) error {
				t.ID = 1
				return nil
			},
		}
		svc := NewService(repo)

		req := CreateRequest{Title: "Test Task", Priority: PriorityHigh}
		todo, err := svc.Create(ctx, req)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if todo.ID != 1 {
			t.Errorf("expected ID 1, got %d", todo.ID)
		}
		if todo.Title != "Test Task" {
			t.Errorf("expected title 'Test Task', got %s", todo.Title)
		}
	})

	t.Run("empty title", func(t *testing.T) {
		svc := NewService(&mockRepository{})
		req := CreateRequest{Title: ""}
		_, err := svc.Create(ctx, req)

		if !errors.Is(err, ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		expectedErr := errors.New("db error")
		repo := &mockRepository{
			createFunc: func(ctx context.Context, t *Todo) error {
				return expectedErr
			},
		}
		svc := NewService(repo)

		req := CreateRequest{Title: "Test"}
		_, err := svc.Create(ctx, req)

		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &Todo{ID: 1, Title: "Test"}
		repo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return expected, nil
			},
		}
		svc := NewService(repo)

		todo, err := svc.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if todo.ID != expected.ID {
			t.Errorf("expected ID %d, got %d", expected.ID, todo.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return nil, ErrNotFound
			},
		}
		svc := NewService(repo)

		_, err := svc.GetByID(ctx, 99)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []Todo{{ID: 1, Title: "T1"}, {ID: 2, Title: "T2"}}
		repo := &mockRepository{
			listFunc: func(ctx context.Context, filter ListFilter) ([]Todo, error) {
				return expected, nil
			},
		}
		svc := NewService(repo)

		todos, err := svc.List(ctx, ListFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(todos) != 2 {
			t.Errorf("expected 2 todos, got %d", len(todos))
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockRepository{
			listFunc: func(ctx context.Context, filter ListFilter) ([]Todo, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewService(repo)

		_, err := svc.List(ctx, ListFilter{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("empty list returns empty slice instead of nil", func(t *testing.T) {
		repo := &mockRepository{
			listFunc: func(ctx context.Context, filter ListFilter) ([]Todo, error) {
				return nil, nil
			},
		}
		svc := NewService(repo)

		todos, err := svc.List(ctx, ListFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if todos == nil {
			t.Error("expected empty slice, got nil")
		}
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success all fields", func(t *testing.T) {
		existing := &Todo{ID: 1, Title: "Old"}
		repo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return existing, nil
			},
			updateFunc: func(ctx context.Context, t *Todo) error {
				return nil
			},
		}
		svc := NewService(repo)

		title := "New"
		desc := "Desc"
		done := true
		priority := PriorityLow
		now := time.Now()
		
		todo, err := svc.Update(ctx, 1, UpdateRequest{
			Title:       &title,
			Description: &desc,
			Done:        &done,
			Priority:    &priority,
			DueDate:     &now,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if todo.Title != title || todo.Description != desc || todo.Done != done || todo.Priority != priority || *todo.DueDate != now {
			t.Error("fields were not updated correctly")
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return nil, ErrNotFound
			},
		}
		svc := NewService(repo)

		_, err := svc.Update(ctx, 1, UpdateRequest{})
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("update error", func(t *testing.T) {
		repo := &mockRepository{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return &Todo{ID: 1}, nil
			},
			updateFunc: func(ctx context.Context, t *Todo) error {
				return errors.New("update failed")
			},
		}
		svc := NewService(repo)

		_, err := svc.Update(ctx, 1, UpdateRequest{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockRepository{
			deleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		svc := NewService(repo)

		err := svc.Delete(ctx, 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
