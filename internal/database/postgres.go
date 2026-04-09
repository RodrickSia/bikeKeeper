// init and setup connection to postgres database
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectToPostgres() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	log.Println("connected to postgres")
	return db, nil
}