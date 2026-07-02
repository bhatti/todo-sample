package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// TodoHandler handles HTTP requests for the /users/{userID}/todos endpoints.
type TodoHandler struct {
	users store.UserStore
	todos store.TodoStore
}

// NewTodoHandler returns a new TodoHandler.
func NewTodoHandler(users store.UserStore, todos store.TodoStore) *TodoHandler {
	return &TodoHandler{users: users, todos: todos}
}

// createTodoRequest is the JSON body for creating a todo.
type createTodoRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	DueAt       *string `json:"due_at"`
}

// updateTodoRequest is the JSON body for updating a todo.
type updateTodoRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	Status      string  `json:"status"`
	DueAt       *string `json:"due_at"`
}

// todoResponse is the JSON representation of a todo.
type todoResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toTodoResponse(t *model.Todo) todoResponse {
	return todoResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    priorityToString(t.Priority),
		Status:      statusToString(t.Status),
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func priorityToString(p model.Priority) string {
	switch p {
	case model.PriorityLow:
		return "low"
	case model.PriorityMedium:
		return "medium"
	case model.PriorityHigh:
		return "high"
	case model.PriorityUrgent:
		return "urgent"
	default:
		return "low"
	}
}

func priorityFromString(s string) (model.Priority, bool) {
	switch strings.ToLower(s) {
	case "low", "":
		return model.PriorityLow, true
	case "medium":
		return model.PriorityMedium, true
	case "high":
		return model.PriorityHigh, true
	case "urgent":
		return model.PriorityUrgent, true
	default:
		return model.PriorityLow, false
	}
}

func statusToString(s model.Status) string {
	switch s {
	case model.StatusPending:
		return "pending"
	case model.StatusInProgress:
		return "in_progress"
	case model.StatusDone:
		return "done"
	case model.StatusCancelled:
		return "cancelled"
	default:
		return "pending"
	}
}

func statusFromString(s string) (model.Status, bool) {
	switch strings.ToLower(s) {
	case "pending", "":
		return model.StatusPending, true
	case "in_progress":
		return model.StatusInProgress, true
	case "done":
		return model.StatusDone, true
	case "cancelled":
		return model.StatusCancelled, true
	default:
		return model.StatusPending, false
	}
}

// CreateTodo handles POST /users/{userID}/todos.
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	// Verify user exists.
	if _, err := h.users.Get(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	priority, ok := priorityFromString(req.Priority)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid priority")
		return
	}

	now := time.Now().UTC()
	todo := &model.Todo{
		ID:          newID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      model.StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if req.DueAt != nil {
		t, err := time.Parse(time.RFC3339, *req.DueAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid due_at format, use RFC3339")
			return
		}
		todo.DueAt = &t
	}

	if err := h.todos.Create(r.Context(), todo); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, toTodoResponse(todo))
}

// GetTodo handles GET /users/{userID}/todos/{id}.
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	id := r.PathValue("id")

	todo, err := h.todos.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if todo.UserID != userID {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, toTodoResponse(todo))
}

// ListTodos handles GET /users/{userID}/todos.
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	todos, err := h.todos.ListByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := make([]todoResponse, 0, len(todos))
	for _, t := range todos {
		resp = append(resp, toTodoResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateTodo handles PUT /users/{userID}/todos/{id}.
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	id := r.PathValue("id")

	existing, err := h.todos.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Priority != "" {
		priority, ok := priorityFromString(req.Priority)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid priority")
			return
		}
		existing.Priority = priority
	}
	if req.Status != "" {
		status, ok := statusFromString(req.Status)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid status")
			return
		}
		existing.Status = status
	}
	if req.DueAt != nil {
		t, err := time.Parse(time.RFC3339, *req.DueAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid due_at format, use RFC3339")
			return
		}
		existing.DueAt = &t
	}
	existing.UpdatedAt = time.Now().UTC()

	if err := h.todos.Update(r.Context(), existing); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, toTodoResponse(existing))
}

// DeleteTodo handles DELETE /users/{userID}/todos/{id}.
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	id := r.PathValue("id")

	existing, err := h.todos.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.todos.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
