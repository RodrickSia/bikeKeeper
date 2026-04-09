package integration_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://bikekeeper:bikekeeper@localhost:5432/bikekeeper?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
