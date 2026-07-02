// Package main is the entry point for the todo HTTP server.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/user/todo/internal/handler"
	"github.com/user/todo/internal/router"
	"github.com/user/todo/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	uh := handler.NewUserHandler(us)
	th := handler.NewTodoHandler(us, ts)
	h := router.New(uh, th)

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
