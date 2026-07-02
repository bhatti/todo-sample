// Package main is the entrypoint for the todo application.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/user/todo/internal/handler"
	"github.com/user/todo/internal/router"
	"github.com/user/todo/internal/store/memory"
)

// main initializes the in-memory store, wires up the HTTP handlers and router,
// and starts the HTTP server on port 8080. It terminates the process with a
// fatal log message if the server fails to start or encounters a fatal error.
func main() {
	// Initialize the in-memory backing store shared by all handlers.
	s := memory.New()

	// Create handlers for user and todo resources, injecting the store.
	userHandler := handler.NewUserHandler(s)
	todoHandler := handler.NewTodoHandler(s, s)

	// Build the request multiplexer with all routes registered.
	mux := router.New(userHandler, todoHandler)

	// Start listening for incoming HTTP connections.
	addr := ":8080"
	fmt.Printf("Starting todo server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
