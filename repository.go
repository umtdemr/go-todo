package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/umtdemr/go-todo/todo"
)

type Repository interface {
	CreateTodo(data *todo.Todo) error
	GetAllTodos() ([]todo.Todo, error)
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
