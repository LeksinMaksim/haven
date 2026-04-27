package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, req CreateRequest) (*Todo, error) {
				return &Todo{ID: 1, Title: req.Title}, nil
			},
		}
		h := NewHandler(svc)

		body := `{"title": "Test Task"}`
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
		rr := httptest.NewRecorder()

		h.create(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rr.Code)
		}

		var resp map[string]Todo
		json.NewDecoder(rr.Body).Decode(&resp)
		if resp["data"].Title != "Test Task" {
			t.Errorf("expected title 'Test Task', got %s", resp["data"].Title)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(`{invalid`))
		rr := httptest.NewRecorder()

		h.create(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, req CreateRequest) (*Todo, error) {
				return nil, ErrEmptyTitle
			},
		}
		h := NewHandler(svc)

		body := `{"title": ""}`
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(body))
		rr := httptest.NewRecorder()

		h.create(rr, req)

		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", rr.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockService{
			createFunc: func(ctx context.Context, req CreateRequest) (*Todo, error) {
				return nil, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("POST", "/api/todos", strings.NewReader(`{"title":"T"}`))
		rr := httptest.NewRecorder()
		h.create(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})
}

func TestHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			listFunc: func(ctx context.Context, filter ListFilter) ([]Todo, error) {
				return []Todo{{ID: 1, Title: "Task 1"}}, nil
			},
		}
		h := NewHandler(svc)

		req := httptest.NewRequest("GET", "/api/todos?done=true&priority=high", nil)
		rr := httptest.NewRecorder()

		h.list(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("invalid done param", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("GET", "/api/todos?done=notabool", nil)
		rr := httptest.NewRecorder()

		h.list(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockService{
			listFunc: func(ctx context.Context, filter ListFilter) ([]Todo, error) {
				return nil, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("GET", "/api/todos", nil)
		rr := httptest.NewRecorder()
		h.list(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return &Todo{ID: id, Title: "Task"}, nil
			},
		}
		h := NewHandler(svc)

		req := httptest.NewRequest("GET", "/api/todos/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		h.getByID(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return nil, ErrNotFound
			},
		}
		h := NewHandler(svc)

		req := httptest.NewRequest("GET", "/api/todos/99", nil)
		req.SetPathValue("id", "99")
		rr := httptest.NewRecorder()

		h.getByID(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("GET", "/api/todos/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		h.getByID(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockService{
			getByIDFunc: func(ctx context.Context, id int64) (*Todo, error) {
				return nil, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("GET", "/api/todos/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.getByID(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, id int64, req UpdateRequest) (*Todo, error) {
				return &Todo{ID: id, Title: *req.Title}, nil
			},
		}
		h := NewHandler(svc)

		newTitle := "Updated"
		body, _ := json.Marshal(UpdateRequest{Title: &newTitle})
		req := httptest.NewRequest("PATCH", "/api/todos/1", bytes.NewReader(body))
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		h.update(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("PATCH", "/api/todos/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()
		h.update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("PATCH", "/api/todos/1", strings.NewReader(`{`))
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, id int64, req UpdateRequest) (*Todo, error) {
				return nil, ErrNotFound
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("PATCH", "/api/todos/1", strings.NewReader(`{}`))
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.update(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockService{
			updateFunc: func(ctx context.Context, id int64, req UpdateRequest) (*Todo, error) {
				return nil, errors.New("db error")
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("PATCH", "/api/todos/1", strings.NewReader(`{}`))
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.update(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		h := NewHandler(svc)

		req := httptest.NewRequest("DELETE", "/api/todos/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		h.delete(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", rr.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		h := NewHandler(&mockService{})
		req := httptest.NewRequest("DELETE", "/api/todos/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()
		h.delete(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int64) error {
				return ErrNotFound
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("DELETE", "/api/todos/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.delete(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		svc := &mockService{
			deleteFunc: func(ctx context.Context, id int64) error {
				return errors.New("db error")
			},
		}
		h := NewHandler(svc)
		req := httptest.NewRequest("DELETE", "/api/todos/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()
		h.delete(rr, req)
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})
}
