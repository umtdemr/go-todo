package todo

import "time"

type Todo struct {
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTodoData struct {
	Title string `json:"title"`
}

func NewTodo(title string) Todo {
	return Todo{
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
