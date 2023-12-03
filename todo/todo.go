package todo

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	Id        uuid.UUID
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTodoData struct {
	Title string `json:"title"`
}

type UpdateTodoData struct {
	Id    string  `json:"id"`
	Title *string `json:"title,omitempty"`
	Done  *bool   `json:"done,omitempty"`
}

func NewTodo(title string) *Todo {
	return &Todo{
		Id:        uuid.New(),
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
