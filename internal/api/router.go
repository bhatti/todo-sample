package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/todo/internal/store"
)

func NewRouter(us store.UserStore, ts store.TodoStore) http.Handler {
	r := chi.NewRouter()

	uh := NewUserHandler(us)
	r.Post("/users", uh.Create)
	r.Get("/users", uh.List)
	r.Get("/users/{id}", uh.GetByID)
	r.Put("/users/{id}", uh.Update)
	r.Delete("/users/{id}", uh.Delete)

	th := NewTodoHandler(us, ts)
	r.Post("/users/{userID}/todos", th.Create)
	r.Get("/users/{userID}/todos", th.List)
	r.Get("/users/{userID}/todos/{id}", th.GetByID)
	r.Put("/users/{userID}/todos/{id}", th.Update)
	r.Delete("/users/{userID}/todos/{id}", th.Delete)

	return r
}
