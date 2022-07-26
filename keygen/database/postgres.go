package database

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func NewPgDb(url string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), url)
}
