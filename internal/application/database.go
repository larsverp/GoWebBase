package application

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase() (*sql.DB, error) {
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
