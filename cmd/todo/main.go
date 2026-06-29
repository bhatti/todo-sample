package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/todo/internal/api"
	"github.com/user/todo/internal/store"
)

func main() {
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8080"
	}

	us := store.NewMemoryUserStore()
	ts := store.NewMemoryTodoStore()
	handler := api.NewRouter(us, ts)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		fmt.Fprintf(os.Stderr, "listening on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}
	fmt.Fprintln(os.Stderr, "server stopped")
}
