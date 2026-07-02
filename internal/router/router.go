// Package router wires HTTP routes to handlers.
package router

import (
	"net/http"

	"github.com/user/todo/internal/handler"
)

// New creates and returns a configured *http.ServeMux with all routes registered.
func New(users *handler.UserHandler, todos *handler.TodoHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", users.Create)
	mux.HandleFunc("GET /users", users.List)
	mux.HandleFunc("GET /users/{userID}", users.Get)
	mux.HandleFunc("DELETE /users/{userID}", users.Delete)

	mux.HandleFunc("POST /users/{userID}/todos", todos.Create)
	mux.HandleFunc("GET /users/{userID}/todos", todos.List)
	mux.HandleFunc("GET /users/{userID}/todos/{todoID}", todos.Get)
	mux.HandleFunc("PUT /users/{userID}/todos/{todoID}", todos.Update)
	mux.HandleFunc("DELETE /users/{userID}/todos/{todoID}", todos.Delete)

	return mux
}
