package files

import (
	"context"
	"os"

	_ "embed"

	"github.com/jackc/pgx/v5"
)

//go:embed schema.sql
var ddl string

func CreateNewDb() (*pgx.Conn, error) {
	ctx := context.Background()
	connString := os.Getenv("POSTGRES_CONNECTION_STRING")
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(ctx, ddl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
