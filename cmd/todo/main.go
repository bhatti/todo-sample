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

func main() {
	s := memory.New()

	userHandler := handler.NewUserHandler(s)
	todoHandler := handler.NewTodoHandler(s, s)

	mux := router.New(userHandler, todoHandler)

	addr := ":8080"
	fmt.Printf("Starting todo server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
