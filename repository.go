package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/umtdemr/go-todo/todo"
)

type Repository interface {
	CreateTodo(data *todo.CreateTodoData) error
}

type PostgresStore struct {
	db *pgx.Conn
}

func (store *PostgresStore) CreateTodo(data *todo.CreateTodoData) error {
	return nil
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
