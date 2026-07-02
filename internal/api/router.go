// Package api provides HTTP handlers and routing for the todo REST API.
package api

import (
	"net/http"
	"strings"

	"github.com/user/todo/internal/store"
)

// NewRouter returns an http.Handler with all API routes registered.
func NewRouter(users store.UserStore, todos store.TodoStore) http.Handler {
	usersH := NewUsersHandler(users)
	todosH := NewTodosHandler(users, todos)

	mux := http.NewServeMux()

	// User routes.
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			usersH.CreateUser(w, r)
		case http.MethodGet:
			usersH.ListUsers(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if path contains /todos to dispatch to todo handler.
		if strings.Contains(path, "/todos") {
			// /users/{id}/todos or /users/{id}/todos/{todoID}
			todoRouteHandler(w, r, todosH)
			return
		}

		// /users/{id}
		switch r.Method {
		case http.MethodGet:
			usersH.GetUser(w, r)
		case http.MethodPut:
			usersH.UpdateUser(w, r)
		case http.MethodDelete:
			usersH.DeleteUser(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	return mux
}

// todoRouteHandler dispatches todo sub-resource requests.
func todoRouteHandler(w http.ResponseWriter, r *http.Request, h *TodosHandler) {
	_, todoID := parseTodosPath(r.URL.Path)
	hasTodoID := todoID != ""

	if hasTodoID {
		// /users/{id}/todos/{todoID}
		switch r.Method {
		case http.MethodGet:
			h.GetTodo(w, r)
		case http.MethodPut:
			h.UpdateTodo(w, r)
		case http.MethodDelete:
			h.DeleteTodo(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	// /users/{id}/todos
	switch r.Method {
	case http.MethodPost:
		h.CreateTodo(w, r)
	case http.MethodGet:
		h.ListTodos(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
