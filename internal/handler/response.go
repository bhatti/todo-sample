package handler

import (
	"time"

	"github.com/user/todo/internal/model"
)

// userResponse is the JSON representation of a user.
type userResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toUserResponse(u model.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserResponses(users []model.User) []userResponse {
	result := make([]userResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(u)
	}
	return result
}

// todoResponse is the JSON representation of a todo.
type todoResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      int        `json:"status"`
	DueAt       *time.Time `json:"due_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toTodoResponse(t model.Todo) todoResponse {
	return todoResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    int(t.Priority),
		Status:      int(t.Status),
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toTodoResponses(todos []model.Todo) []todoResponse {
	result := make([]todoResponse, len(todos))
	for i, t := range todos {
		result[i] = toTodoResponse(t)
	}
	return result
}
