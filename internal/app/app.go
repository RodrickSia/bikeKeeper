package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
	"github.com/RodrickSia/bikeKeeper/internal/card"
	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
	"github.com/RodrickSia/bikeKeeper/internal/payment"
	"github.com/RodrickSia/bikeKeeper/internal/user"
	"github.com/RodrickSia/bikeKeeper/pkg/OCR"
	"github.com/RodrickSia/bikeKeeper/pkg/storage"
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

func (a *App) registerRoutes(prefix string) {
	// auth + users (shared repo)
	userRepo := user.NewRepository(a.DB)
	authAdapter := user.NewAuthAdapter(userRepo)
	authSvc := auth.NewService(authAdapter)
	authHandler := auth.NewHandler(authSvc)
	auth.RegisterRoutes(a.Router, authHandler, prefix)

	// middleware chains
	authenticated := auth.Authenticate(authSvc)
	staffOnly := auth.RequireRole(user.RoleStaff, user.RoleFaculty, user.RoleAdmin)
	facultyOnly := auth.RequireRole(user.RoleFaculty, user.RoleAdmin)

	// payment
	paymentRepo := payment.NewRepository(a.DB)
	paymentSvc := payment.NewService(paymentRepo)
	paymentHandler := payment.NewHandler(paymentSvc)
	payment.RegisterRoutes(a.Router, paymentHandler, prefix, authenticated)

	// users
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)
	user.RegisterRoutes(a.Router, userHandler, prefix, authenticated, facultyOnly)

	// parking sessions
	sessionRepo := parkingsession.NewRepository(a.DB)
	imageStore := storage.NewLocalStorage("./images")
	plateProcessor := OCR.NewPlateProcessor()
	sessionSvc := parkingsession.NewService(sessionRepo, plateProcessor, imageStore, &paymentAdapter{svc: paymentSvc})
	sessionHandler := parkingsession.NewHandler(sessionSvc)
	parkingsession.RegisterRoutes(a.Router, sessionHandler, prefix, authenticated, staffOnly)

	// members
	memberRepo := member.NewRepository(a.DB)
	memberSvc := member.NewService(memberRepo)
	memberHandler := member.NewHandler(memberSvc)
	member.RegisterRoutes(a.Router, memberHandler, prefix, authenticated, facultyOnly)

	// cards
	cardRepo := card.NewRepository(a.DB)
	cardSvc := card.NewService(cardRepo)
	cardHandler := card.NewHandler(cardSvc)
	card.RegisterRoutes(a.Router, cardHandler, prefix, authenticated, facultyOnly)
}

type paymentAdapter struct {
	svc *payment.Service
}

func (a *paymentAdapter) ChargeParking(ctx context.Context, cardUID string, fee float64, sessionID int64) error {
	_, err := a.svc.ChargeParking(ctx, cardUID, fee, sessionID)
	return err
}
