// Package router wires HTTP routes to handlers.
package router

import (
	"net/http"

	"github.com/user/todo/internal/handler"
)

// New returns an http.Handler with all routes registered.
func New(uh *handler.UserHandler, th *handler.TodoHandler) http.Handler {
	mux := http.NewServeMux()

	// User routes.
	mux.HandleFunc("POST /users", uh.CreateUser)
	mux.HandleFunc("GET /users", uh.ListUsers)
	mux.HandleFunc("GET /users/{id}", uh.GetUser)
	mux.HandleFunc("PUT /users/{id}", uh.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", uh.DeleteUser)

	// Todo routes (nested under users).
	mux.HandleFunc("POST /users/{userID}/todos", th.CreateTodo)
	mux.HandleFunc("GET /users/{userID}/todos", th.ListTodos)
	mux.HandleFunc("GET /users/{userID}/todos/{id}", th.GetTodo)
	mux.HandleFunc("PUT /users/{userID}/todos/{id}", th.UpdateTodo)
	mux.HandleFunc("DELETE /users/{userID}/todos/{id}", th.DeleteTodo)

	return mux
}
