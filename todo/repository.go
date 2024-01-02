package todo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

type IRepository interface {
	CreateTodo(data *Todo, userId int64) (*Todo, error)
	GetAllTodos(userId int64) ([]Todo, error)
	GetTodo(todoId int, userId int64) (*Todo, error)
	UpdateTodo(data *UpdateTodoData, userId int64) (*Todo, error)
	RemoveTodo(todoId int, userId int64) (*Todo, error)
}

type Repository struct {
	DB *pgx.Conn
}

func NewTodoRepository(dbConn *pgx.Conn) (*Repository, error) {
	return &Repository{dbConn}, nil
}

func (store *Repository) Init() error {
	return store.CreateTodoTable()
}

func (store *Repository) CreateTodoTable() error {
	query := `CREATE TABLE IF NOT EXISTS "todo" (
		id serial PRIMARY KEY,
		user_id integer NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
		title varchar(255) NOT NULL,
		done boolean DEFAULT false,
		created_at timestamp DEFAULT now(),
		updated_at timestamp DEFAULT now()
	)`

	_, err := store.DB.Exec(context.Background(), query)
	return err
}

func (store *Repository) CreateTodo(data *Todo, userId int64) (*Todo, error) {
	query := `INSERT INTO "todo"(title, user_id) VALUES (@title, @userId) RETURNING id, title, done, created_at, updated_at`
	args := pgx.NamedArgs{
		"title":  data.Title,
		"userId": userId,
	}

	rows := store.DB.QueryRow(context.Background(), query, args)
	createdTodo, scanErr := ScanTodo(rows)
	if scanErr != nil {
		return nil, scanErr
	}

	return createdTodo, nil
}

func (store *Repository) GetAllTodos(userId int64) ([]Todo, error) {
	query := `SELECT id, title, done, created_at, updated_at FROM "todo" WHERE user_id=@userId`
	args := pgx.NamedArgs{"userId": userId}

	rows, err := store.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var t Todo

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

func (store *Repository) UpdateTodo(data *UpdateTodoData, userId int64) (*Todo, error) {
	var updateBuilder strings.Builder
	updateBuilder.WriteString("UPDATE todo SET ")
	var updates []string
	var args []interface{}

	if data.Title == nil && data.Done == nil {
		return nil, errors.New("no field is provided")
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

	updateBuilder.WriteString(fmt.Sprintf("WHERE id = %d and user_id = %d RETURNING id, title, done, created_at, updated_at", *data.Id, userId))

	rows := store.DB.QueryRow(context.Background(), updateBuilder.String(), args...)

	updatedData, updateScanErr := ScanTodo(rows)

	if updateScanErr != nil {
		return nil, updateScanErr
	}

	return updatedData, nil
}

func (store *Repository) GetTodo(todoId int, userId int64) (*Todo, error) {
	query := `SELECT id, title, done, created_at, updated_at FROM todo WHERE id = @todoId and user_id = @userId`
	args := pgx.NamedArgs{
		"todoId": todoId,
		"userId": userId,
	}

	var singleTodo *Todo
	singleTodo = new(Todo)
	row := store.DB.QueryRow(context.Background(), query, args)

	err := row.Scan(&singleTodo.Id, &singleTodo.Title, &singleTodo.Done, &singleTodo.CreatedAt, &singleTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return singleTodo, nil
}

func (store *Repository) RemoveTodo(todoId int, userId int64) (*Todo, error) {
	query := `DELETE FROM todo WHERE id = @todoId and user_id = @userId RETURNING id, title, done, created_at, updated_at`

	args := pgx.NamedArgs{
		"todoId": todoId,
		"userId": userId,
	}

	rows := store.DB.QueryRow(context.Background(), query, args)
	removedTodo, removeScanErr := ScanTodo(rows)

	if removeScanErr != nil {
		return nil, removeScanErr
	}

	return removedTodo, nil
}
