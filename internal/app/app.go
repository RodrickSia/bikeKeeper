package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
	"github.com/RodrickSia/bikeKeeper/internal/card"
	"github.com/RodrickSia/bikeKeeper/internal/cardrequest"
	"github.com/RodrickSia/bikeKeeper/internal/device"
	"github.com/RodrickSia/bikeKeeper/internal/incident"
	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/internal/monthlypass"
	"github.com/RodrickSia/bikeKeeper/internal/notification"
	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
	"github.com/RodrickSia/bikeKeeper/internal/parkinglot"
	"github.com/RodrickSia/bikeKeeper/internal/payment"
	"github.com/RodrickSia/bikeKeeper/internal/shift"
	"github.com/RodrickSia/bikeKeeper/internal/support"
	"github.com/RodrickSia/bikeKeeper/internal/user"
	"github.com/RodrickSia/bikeKeeper/internal/vehicle"
	"github.com/RodrickSia/bikeKeeper/internal/visitor"
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
	// Health check — no auth required
	a.Router.HandleFunc("GET "+prefix+"/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// core repos + services initialized early for auth registration wiring
	userRepo := user.NewRepository(a.DB)
	userSvc := user.NewService(userRepo)

	memberRepo := member.NewRepository(a.DB)
	memberSvc := member.NewService(memberRepo)

	authAdapter := user.NewAuthAdapter(userRepo)
	authSvc := auth.NewService(authAdapter)
	authHandler := auth.NewHandler(authSvc, &userCreatorAdapter{svc: userSvc}, &memberCreatorAdapter{svc: memberSvc})
	auth.RegisterRoutes(a.Router, authHandler, prefix)

	// middleware chains
	authenticated := auth.Authenticate(authSvc)
	staffOnly := auth.RequireRole(user.RoleStaff, user.RoleFaculty, user.RoleAdmin)
	facultyOnly := auth.RequireRole(user.RoleFaculty, user.RoleAdmin)

	// cards
	cardRepo := card.NewRepository(a.DB)
	cardSvc := card.NewService(cardRepo)
	cardHandler := card.NewHandler(cardSvc)
	card.RegisterRoutes(a.Router, cardHandler, prefix, authenticated)

	// payment
	paymentRepo := payment.NewRepository(a.DB)
	paymentSvc := payment.NewService(paymentRepo)
	paymentHandler := payment.NewHandler(paymentSvc, &cardFinderAdapter{repo: cardRepo})
	payment.RegisterRoutes(a.Router, paymentHandler, prefix, authenticated)

	// users
	userHandler := user.NewHandler(userSvc)
	user.RegisterRoutes(a.Router, userHandler, prefix, authenticated, facultyOnly)

	// parking sessions
	sessionRepo := parkingsession.NewRepository(a.DB)
	imageStore := storage.NewLocalStorage("./images")
	plateProcessor, err := OCR.NewPlateProcessor()
	if err != nil {
		panic(err)
	}
	sessionSvc := parkingsession.NewService(sessionRepo, plateProcessor, imageStore, &paymentAdapter{svc: paymentSvc})
	sessionHandler := parkingsession.NewHandler(sessionSvc)
	parkingsession.RegisterRoutes(a.Router, sessionHandler, prefix, authenticated, staffOnly)

	// members
	memberHandler := member.NewHandler(memberSvc)
	member.RegisterRoutes(a.Router, memberHandler, prefix, authenticated, facultyOnly)

	// monthly passes
	monthlyPassRepo := monthlypass.NewRepository(a.DB)
	monthlyPassSvc := monthlypass.NewService(monthlyPassRepo)
	monthlyPassHandler := monthlypass.NewHandler(monthlyPassSvc)
	monthlypass.RegisterRoutes(a.Router, monthlyPassHandler, prefix, authenticated)

	// card requests
	cardRequestRepo := cardrequest.NewRepository(a.DB)
	cardRequestSvc := cardrequest.NewService(cardRequestRepo, &cardRepoAdapter{repo: cardRepo})
	cardRequestHandler := cardrequest.NewHandler(cardRequestSvc)
	cardrequest.RegisterRoutes(a.Router, cardRequestHandler, prefix, authenticated, facultyOnly)

	// parking lots
	parkinglotRepo := parkinglot.NewRepository(a.DB)
	parkinglotSvc := parkinglot.NewService(parkinglotRepo)
	parkinglotHandler := parkinglot.NewHandler(parkinglotSvc)
	parkinglot.RegisterRoutes(a.Router, parkinglotHandler, prefix, authenticated, facultyOnly)

	// shifts
	shiftRepo := shift.NewRepository(a.DB)
	shiftSvc := shift.NewService(shiftRepo)
	shiftHandler := shift.NewHandler(shiftSvc)
	shift.RegisterRoutes(a.Router, shiftHandler, prefix, authenticated, staffOnly)

	// incidents
	incidentRepo := incident.NewRepository(a.DB)
	incidentSvc := incident.NewService(incidentRepo)
	incidentHandler := incident.NewHandler(incidentSvc)
	incident.RegisterRoutes(a.Router, incidentHandler, prefix, authenticated, staffOnly)

	// devices
	deviceRepo := device.NewRepository(a.DB)
	deviceSvc := device.NewService(deviceRepo)
	deviceHandler := device.NewHandler(deviceSvc)
	device.RegisterRoutes(a.Router, deviceHandler, prefix, authenticated, staffOnly)

	// notifications
	notificationRepo := notification.NewRepository(a.DB)
	notificationSvc := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationSvc)
	notification.RegisterRoutes(a.Router, notificationHandler, prefix, authenticated, staffOnly)

	// support tickets
	supportRepo := support.NewRepository(a.DB)
	supportSvc := support.NewService(supportRepo)
	supportHandler := support.NewHandler(supportSvc)
	support.RegisterRoutes(a.Router, supportHandler, prefix, authenticated, staffOnly)

	// visitor passes
	visitorRepo := visitor.NewRepository(a.DB)
	visitorSvc := visitor.NewService(visitorRepo)
	visitorHandler := visitor.NewHandler(visitorSvc)
	visitor.RegisterRoutes(a.Router, visitorHandler, prefix, authenticated)

	// vehicles
	vehicleRepo := vehicle.NewRepository(a.DB)
	vehicleSvc := vehicle.NewService(vehicleRepo)
	vehicleHandler := vehicle.NewHandler(vehicleSvc)
	vehicle.RegisterRoutes(a.Router, vehicleHandler, prefix, authenticated, staffOnly)
}

type paymentAdapter struct {
	svc *payment.Service
}

func (a *paymentAdapter) ChargeParking(ctx context.Context, cardUID string, fee float64, sessionID int64) error {
	_, err := a.svc.ChargeParking(ctx, cardUID, fee, sessionID)
	return err
}

type cardRepoAdapter struct {
	repo card.Repository
}

func (a *cardRepoAdapter) CreateCard(ctx context.Context, cardUID, cardType string, memberID *string) error {
	return a.repo.Create(ctx, &card.Card{
		CardUID:  cardUID,
		CardType: cardType,
		MemberID: memberID,
		IsInside: false,
		Status:   "active",
		Balance:  0.0,
	})
}

type userCreatorAdapter struct {
	svc *user.Service
}

func (a *userCreatorAdapter) CreateRegister(ctx context.Context, email, password, role string, memberID *string) error {
	_, err := a.svc.Create(ctx, user.CreateParams{
		Email:    email,
		Password: password,
		Role:     role,
		MemberID: memberID,
		Status:   user.StatusPending,
	})
	return err
}

type cardFinderAdapter struct {
	repo card.Repository
}

func (a *cardFinderAdapter) GetByUID(ctx context.Context, cardUID string) (*payment.CardInfo, error) {
	c, err := a.repo.GetByUID(ctx, cardUID)
	if err != nil {
		return nil, err
	}
	return &payment.CardInfo{CardUID: c.CardUID, MemberID: c.MemberID}, nil
}

type memberCreatorAdapter struct {
	svc *member.Service
}

func (a *memberCreatorAdapter) CreateRegister(ctx context.Context, studentID, fullName string, phone *string) (string, error) {
	m, err := a.svc.Create(ctx, member.CreateParams{
		StudentID: studentID,
		FullName:  fullName,
		Phone:     phone,
	})
	if err != nil {
		return "", err
	}
	return m.ID, nil
}
