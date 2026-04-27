package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/LeksinMaksim/haven/pkg/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/todos", h.list)
	mux.HandleFunc("POST /api/todos", h.create)
	mux.HandleFunc("GET /api/todos/{id}", h.getByID)
	mux.HandleFunc("PATCH /api/todos/{id}", h.update)
	mux.HandleFunc("DELETE /api/todos/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	filter := ListFilter{}

	if d := r.URL.Query().Get("done"); d != "" {
		done, err := strconv.ParseBool(d)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid 'done' param")
			return
		}
		filter.Done = &done
	}

	if p := r.URL.Query().Get("priority"); p != "" {
		priority := Priority(p)
		filter.Priority = &priority
	}

	todos, err := h.svc.List(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list todos")
		return
	}

	response.OK(w, todos)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.svc.Create(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmptyTitle) {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create todo")
		return
	}

	response.Created(w, todo)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get todo")
		return
	}

	response.OK(w, todo)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update todo")
		return
	}

	response.OK(w, todo)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
