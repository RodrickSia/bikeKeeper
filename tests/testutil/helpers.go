package testutil

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://bikekeeper:bikekeeper@localhost:5432/bikekeeper?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func CleanTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			t.Fatalf("failed to truncate %s: %v", table, err)
		}
	}
}

func StringPtr(s string) *string {
	return &s
}
