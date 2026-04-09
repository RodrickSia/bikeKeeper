package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/card"
	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
	"github.com/RodrickSia/bikeKeeper/tests/testutil"
	"github.com/shopspring/decimal"
)

func setupSessionService(t *testing.T) (*parkingsession.Service, *card.Service, *member.Service, func()) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	sessionRepo := parkingsession.NewRepository(db)
	sessionSvc := parkingsession.NewService(sessionRepo)
	cardRepo := card.NewRepository(db)
	cardSvc := card.NewService(cardRepo)
	memberRepo := member.NewRepository(db)
	memberSvc := member.NewService(memberRepo)
	cleanup := func() {
		testutil.CleanTables(t, db, "parking_sessions", "cards", "members")
	}
	cleanup()
	return sessionSvc, cardSvc, memberSvc, cleanup
}

func seedCard(t *testing.T, cardSvc *card.Service, uid string) {
	t.Helper()
	_, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  uid,
		CardType: "casual",
	})
	if err != nil {
		t.Fatalf("failed to seed card: %v", err)
	}
}

func TestSessionCheckIn(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-001")

	plateIn := "59F1-12345"
	session, err := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: "CARD-001",
		PlateIn: &plateIn,
	})
	if err != nil {
		t.Fatalf("CheckIn failed: %v", err)
	}
	if session.ID == 0 {
		t.Error("expected non-zero session ID")
	}
	if session.CardUID != "CARD-001" {
		t.Errorf("expected CardUID 'CARD-001', got '%s'", session.CardUID)
	}
	if session.Status != "ongoing" {
		t.Errorf("expected Status 'ongoing', got '%s'", session.Status)
	}
	if session.CheckInTime.IsZero() {
		t.Error("expected non-zero CheckInTime")
	}
}

func TestSessionCheckIn_DuplicateBlocked(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-DUP")

	_, err := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-DUP"})
	if err != nil {
		t.Fatalf("first CheckIn failed: %v", err)
	}

	_, err = sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-DUP"})
	if err == nil {
		t.Error("expected error on duplicate check-in, got nil")
	}
}

func TestSessionCheckOut(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-OUT")

	session, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: "CARD-OUT",
		PlateIn: testutil.StringPtr("59F1-99999"),
	})

	plateOut := "59F1-99999"
	err := sessionSvc.CheckOut(ctx, session.ID, parkingsession.CheckOutParams{
		PlateOut:  &plateOut,
		Cost:      decimal.NewFromInt(5000),
		IsWarning: false,
	})
	if err != nil {
		t.Fatalf("CheckOut failed: %v", err)
	}
}

func TestSessionCheckOut_AlreadyCompleted(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-DONE")

	session, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-DONE"})

	sessionSvc.CheckOut(ctx, session.ID, parkingsession.CheckOutParams{
		Cost: decimal.NewFromInt(0),
	})

	err := sessionSvc.CheckOut(ctx, session.ID, parkingsession.CheckOutParams{
		Cost: decimal.NewFromInt(0),
	})
	if err == nil {
		t.Error("expected error on double checkout, got nil")
	}
}

func TestSessionCheckOut_PlateWarning(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-WARN")

	session, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: "CARD-WARN",
		PlateIn: testutil.StringPtr("59F1-11111"),
	})

	plateOut := "59F1-22222"
	err := sessionSvc.CheckOut(ctx, session.ID, parkingsession.CheckOutParams{
		PlateOut:  &plateOut,
		Cost:      decimal.NewFromInt(5000),
		IsWarning: true,
	})
	if err != nil {
		t.Fatalf("CheckOut with warning failed: %v", err)
	}
}

func TestSessionGetByID(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-GET")

	created, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-GET"})

	got, err := sessionSvc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.CardUID != "CARD-GET" {
		t.Errorf("expected CardUID 'CARD-GET', got '%s'", got.CardUID)
	}
}

func TestSessionGetByID_NotFound(t *testing.T) {
	sessionSvc, _, _, cleanup := setupSessionService(t)
	defer cleanup()

	_, err := sessionSvc.GetByID(context.Background(), 999999)
	if err == nil {
		t.Error("expected error for non-existent session, got nil")
	}
}

func TestSessionListByCard(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-LIST")

	// create and complete a session
	s1, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-LIST"})
	sessionSvc.CheckOut(ctx, s1.ID, parkingsession.CheckOutParams{Cost: decimal.NewFromInt(0)})

	// create another ongoing session
	sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-LIST"})

	sessions, err := sessionSvc.ListByCard(ctx, "CARD-LIST")
	if err != nil {
		t.Fatalf("ListByCard failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestSessionDelete(t *testing.T) {
	sessionSvc, cardSvc, _, cleanup := setupSessionService(t)
	defer cleanup()

	ctx := context.Background()
	seedCard(t, cardSvc, "CARD-DEL")

	session, _ := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{CardUID: "CARD-DEL"})

	err := sessionSvc.Delete(ctx, session.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = sessionSvc.GetByID(ctx, session.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestSessionDelete_NotFound(t *testing.T) {
	sessionSvc, _, _, cleanup := setupSessionService(t)
	defer cleanup()

	err := sessionSvc.Delete(context.Background(), 999999)
	if err == nil {
		t.Error("expected error for non-existent session, got nil")
	}
}
