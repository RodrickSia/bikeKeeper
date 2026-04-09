package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embededMigration embed.FS

func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(embededMigration)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, "migrations")
}
func RollbackMigrations(db *sql.DB) error {
	goose.SetBaseFS(embededMigration)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Down(db, "migrations")
}