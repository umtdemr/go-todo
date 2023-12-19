package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/umtdemr/go-todo/todo"
	"strings"
	"time"
)

type Repository interface {
	CreateTodo(data *todo.Todo) error
	GetAllTodos() ([]todo.Todo, error)
	GetTodo(todoId int) (*todo.Todo, error)
	UpdateTodo(data *todo.UpdateTodoData) error
	RemoveTodo(todoId int) error
}

type PostgresStore struct {
	db *pgx.Conn
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{conn}, nil
}

func (store *PostgresStore) Init() error {
	return store.CreateTodoTable()
}

func (store *PostgresStore) CreateTodoTable() error {
	query := `CREATE TABLE IF NOT EXISTS todo(
		id serial PRIMARY KEY,
		title varchar(255),
		done boolean DEFAULT false,
		created_at timestamp DEFAULT now(),
		updated_at timestamp DEFAULT now()
	)`

	_, err := store.db.Exec(context.Background(), query)
	return err
}

func (store *PostgresStore) CreateTodo(data *todo.Todo) error {
	query := `INSERT INTO todo(title) VALUES (@title)`
	args := pgx.NamedArgs{
		"title": data.Title,
	}

	_, err := store.db.Exec(context.Background(), query, args)

	return err
}

func (store *PostgresStore) GetAllTodos() ([]todo.Todo, error) {
	query := `SELECT * FROM todo`

	rows, err := store.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []todo.Todo

	for rows.Next() {
		var t todo.Todo

		err := rows.Scan(&t.Id, &t.Title, &t.Done, &t.CreatedAt, &t.UpdatedAt)

		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (store *PostgresStore) UpdateTodo(data *todo.UpdateTodoData) error {
	var updateBuilder strings.Builder
	updateBuilder.WriteString("UPDATE todo SET ")
	var updates []string
	var args []interface{}

	if data.Title == nil && data.Done == nil {
		return errors.New("no field is provided")
	}

	if data.Title != nil {
		updates = append(updates, "title = $1")
		args = append(args, data.Title)
	}

	if data.Done != nil {
		updates = append(updates, fmt.Sprintf("done = $%d", len(args)+1))
		args = append(args, data.Done)
	}

	updateBuilder.WriteString(strings.Join(updates, ", "))
	updateBuilder.WriteString(",") // Add space before update

	updateBuilder.WriteString(fmt.Sprintf("updated_at = $%d", len(args)+1))
	args = append(args, time.Now())
	updateBuilder.WriteString(" ") // Add space before WHERE clause

	updateBuilder.WriteString(fmt.Sprintf("WHERE id = %d", *data.Id))

	updateResponse, err := store.db.Exec(context.Background(), updateBuilder.String(), args...)

	if err != nil {
		return err
	}

	if updateResponse.RowsAffected() == 0 {
		return errors.New("couldn't update")
	}

	return nil
}

func (store *PostgresStore) GetTodo(todoId int) (*todo.Todo, error) {
	query := fmt.Sprintf(`SELECT * FROM todo WHERE id = %d`, todoId)

	var singleTodo *todo.Todo
	singleTodo = new(todo.Todo)
	row := store.db.QueryRow(context.Background(), query)

	err := row.Scan(&singleTodo.Id, &singleTodo.Title, &singleTodo.Done, &singleTodo.CreatedAt, &singleTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return singleTodo, nil
}

func (store *PostgresStore) RemoveTodo(todoId int) error {
	query := `DELETE FROM todo WHERE id = @todoId`

	args := pgx.NamedArgs{
		"todoId": todoId,
	}

	deleteResponse, err := store.db.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	if deleteResponse.RowsAffected() == 0 {
		return errors.New("couldn't delete")
	}

	return nil
}
