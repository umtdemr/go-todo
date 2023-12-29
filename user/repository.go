package user

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type IRepository interface {
	CreateUser(data *CreateUserData) error
}

type Repository struct {
	db *pgx.Conn
}

func NewUserRepository(dbConn *pgx.Conn) (*Repository, error) {
	return &Repository{dbConn}, nil
}
func (repository *Repository) Init() error {
	return repository.CreateUserTable()
}

func (repository *Repository) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS "user" (
		id serial PRIMARY KEY,
		username varchar(20) NOT NULL,
		password text NOT NULL,
		email varchar(255) NOT NULL,
		created_at timestamp DEFAULT now(),
		is_active bool DEFAULT true,
		is_verified bool DEFAULT false
	)`

	_, err := repository.db.Exec(context.Background(), query)
	return err
}

func (repository *Repository) CreateUser(data *CreateUserData) error {
	query := `INSERT INTO "user"(username, password, email) VALUES (@username, @password, @email)`
	args := pgx.NamedArgs{
		"username": data.Username,
		"password": data.Password,
		"email":    data.Email,
	}

	_, err := repository.db.Exec(context.Background(), query, args)
	return err
}
