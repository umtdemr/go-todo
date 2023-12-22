package main

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type PostgresStore struct {
	DB *pgx.Conn
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{conn}, nil
}
