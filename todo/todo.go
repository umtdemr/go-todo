package todo

import (
	"github.com/jackc/pgx/v5"
	"time"
)

type Todo struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTodoData struct {
	Title string `json:"title"`
}

type UpdateTodoData struct {
	Id    *int    `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
	Done  *bool   `json:"done,omitempty"`
}

type DeleteTodoData struct {
	Id *string `json:"id,omitempty"`
}

func ScanTodo(row pgx.Row) (*Todo, error) {
	var t *Todo
	t = new(Todo) // initialize it since we need to pass values into a pointer
	err := row.Scan(&t.Id, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewTodo(title string) *Todo {
	return &Todo{
		Id:        1,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
