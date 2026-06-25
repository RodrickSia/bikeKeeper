package main

import (
	"log"
	"net/http"
	"os"

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
		log.Printf("Failed to run database migrations: %v", err)
		if err := database.RollbackMigrations(db); err != nil {
			log.Printf("Failed to rollback migrations: %v", err)
		}
		os.Exit(1)
	}
	application := app.New(db, prefix)
	
	log.Println("Initiating Server")
	server := http.Server{
		Addr:    ":8080",
		Handler: enableCORS(application.Router),
	}
	server.ListenAndServe()
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}	

