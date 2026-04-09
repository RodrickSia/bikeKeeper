package app

import (
	"database/sql"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
)

type App struct {
	Router *http.ServeMux
	DB     *sql.DB
}

func New(db *sql.DB, prefix string) *App {
	router := http.NewServeMux()

	a := &App{
		Router: router,
		DB:     db,
	}

	a.registerRoutes(prefix)
	return a
}

// This bind the registreRoutes to the app struct
func (a *App) registerRoutes(prefix string) {
	// parking sessions
	sessionRepo := parkingsession.NewRepository(a.DB)
	sessionSvc := parkingsession.NewService(sessionRepo)
	sessionHandler := parkingsession.NewHandler(sessionSvc)
	parkingsession.RegisterRoutes(a.Router, sessionHandler, prefix)
}
