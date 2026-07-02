package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/user/todo/internal/model"
	"github.com/user/todo/internal/store"
)

// TodosHandler handles HTTP requests for todo resources scoped to a user.
type TodosHandler struct {
	users store.UserStore
	todos store.TodoStore
}

// NewTodosHandler returns a TodosHandler backed by the given stores.
func NewTodosHandler(users store.UserStore, todos store.TodoStore) *TodosHandler {
	return &TodosHandler{users: users, todos: todos}
}

// createTodoRequest is the request body for POST /users/{id}/todos.
type createTodoRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    model.Priority `json:"priority"`
	DueAt       *time.Time     `json:"due_at"`
}

// updateTodoRequest is the request body for PUT /users/{id}/todos/{todoID}.
type updateTodoRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    model.Priority `json:"priority"`
	Status      model.Status   `json:"status"`
	DueAt       *time.Time     `json:"due_at"`
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
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

func toTodoResponse(t *model.Todo) todoResponse {
	return todoResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority.String(),
		Status:      t.Status.String(),
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// parseTodosPath extracts userID and optional todoID from a path of the form:
//
//	/users/{userID}/todos
//	/users/{userID}/todos/{todoID}
func parseTodosPath(path string) (userID, todoID string) {
	// Strip leading /users/
	rest := strings.TrimPrefix(path, "/users/")
	// rest is now "{userID}/todos" or "{userID}/todos/{todoID}"
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) >= 1 {
		userID = parts[0]
	}
	// parts[1] should be "todos"
	if len(parts) >= 3 {
		todoID = parts[2]
	}
	return userID, todoID
}

// CreateTodo handles POST /users/{id}/todos.
func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := parseTodosPath(r.URL.Path)
	if _, err := h.users.GetByID(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req createTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := h.todos.Create(&model.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueAt:       req.DueAt,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toTodoResponse(todo))
}

// ListTodos handles GET /users/{id}/todos.
func (h *TodosHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	userID, _ := parseTodosPath(r.URL.Path)
	if _, err := h.users.GetByID(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	todos, err := h.todos.ListByUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]todoResponse, 0, len(todos))
	for _, t := range todos {
		resp = append(resp, toTodoResponse(t))
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetTodo handles GET /users/{id}/todos/{todoID}.
func (h *TodosHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID, todoID := parseTodosPath(r.URL.Path)
	if _, err := h.users.GetByID(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	todo, err := h.todos.GetByID(todoID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if todo.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found for this user")
		return
	}
	writeJSON(w, http.StatusOK, toTodoResponse(todo))
}

// UpdateTodo handles PUT /users/{id}/todos/{todoID}.
func (h *TodosHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, todoID := parseTodosPath(r.URL.Path)
	if _, err := h.users.GetByID(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existing, err := h.todos.GetByID(todoID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found for this user")
		return
	}

	var req updateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	updated, err := h.todos.Update(&model.Todo{
		ID:          todoID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DueAt:       req.DueAt,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toTodoResponse(updated))
}

// DeleteTodo handles DELETE /users/{id}/todos/{todoID}.
func (h *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, todoID := parseTodosPath(r.URL.Path)
	if _, err := h.users.GetByID(userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	existing, err := h.todos.GetByID(todoID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if existing.UserID != userID {
		writeError(w, http.StatusNotFound, "todo not found for this user")
		return
	}

	if err := h.todos.Delete(todoID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
