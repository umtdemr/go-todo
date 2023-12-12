package todo

import (
	"time"
)

type Todo struct {
	Id        int
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

type DeleteTodoData struct {
	Id *string `json:"id,omitempty"`
}

func NewTodo(title string) *Todo {
	return &Todo{
		Id:        1,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
