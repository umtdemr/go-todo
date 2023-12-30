package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type IRepository interface {
	CreateUser(data *CreateUserData) error
	GetUserWithAllParams(data *LoginUserData) (*UserParams, error)
	GetUserByUsername(username string) *VisibleUser
	GetUserByEmail(email string) *VisibleUser
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

func (repository *Repository) GetUserWithAllParams(data *LoginUserData) (*UserParams, error) {
	var credentialColumnValue, credentialColumnName string
	if data.Username != nil {
		credentialColumnValue = *data.Username
		credentialColumnName = "username"
	} else {
		credentialColumnValue = *data.Email
		credentialColumnName = "email"
	}

	query := fmt.Sprintf(
		`SELECT id, username, email, password, is_active, is_verified, created_at FROM "user" WHERE %v=@credentialVal`,
		credentialColumnName,
	)
	args := pgx.NamedArgs{
		"credentialVal": credentialColumnValue,
	}

	var user UserParams
	queryRow := repository.db.QueryRow(context.Background(), query, args)

	err := queryRow.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.IsVerified,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *Repository) GetUserByUsername(username string) *VisibleUser {
	var user VisibleUser

	query := `SELECT id, username, email, created_at FROM "user" WHERE username=@username`
	args := pgx.NamedArgs{"username": username}

	queryRow := repository.db.QueryRow(context.Background(), query, args)

	err := queryRow.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil
	}

	return &user
}

func (repository *Repository) GetUserByEmail(email string) *VisibleUser {
	var user VisibleUser

	query := `SELECT id, username, email, created_at FROM "user" WHERE email=@email`
	args := pgx.NamedArgs{"email": email}

	queryRow := repository.db.QueryRow(context.Background(), query, args)

	err := queryRow.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil
	}

	return &user
}
