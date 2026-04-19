package main

import (
	"log"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/app"
	"github.com/RodrickSia/bikeKeeper/internal/database"
)

var prefix = "/api/v1"

func main() {
	db, err := database.ConnectToPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()
	// run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
		database.RollbackMigrations(db)
		return
	}
	application := app.New(db, prefix)
	
	log.Println("Initiating Server")
	server := http.Server{
		Addr:    ":8080",
		Handler: application.Router,
	}
	server.ListenAndServe()
}	

