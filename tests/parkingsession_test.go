package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
)

// --- CheckIn ---

func TestSessionCheckIn(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardActive has a *completed* session — no ongoing, so new CheckIn is allowed
	session, err := sessionSvc.CheckIn(context.Background(), parkingsession.CheckInParams{
		CardUID: f.CardActive.CardUID,
	})
	if err != nil {
		t.Fatalf("CheckIn failed: %v", err)
	}
	if session.ID == 0 {
		t.Error("expected non-zero session ID")
	}
	if session.CardUID != f.CardActive.CardUID {
		t.Errorf("CardUID: got %q, want %q", session.CardUID, f.CardActive.CardUID)
	}
	if session.Status != "ongoing" {
		t.Errorf("Status: got %q, want %q", session.Status, "ongoing")
	}
	if session.CheckInTime.IsZero() {
		t.Error("expected non-zero CheckInTime")
	}
}

func TestSessionCheckIn_DuplicateBlocked(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardCasual already has an ongoing session (SessionOngoing in fixtures)
	_, err := sessionSvc.CheckIn(context.Background(), parkingsession.CheckInParams{
		CardUID: f.CardCasual.CardUID,
	})
	if err == nil {
		t.Error("expected error: card already has an ongoing session")
	}
}

// --- CheckOut ---

func TestSessionCheckOut(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// Check out the ongoing session on CardCasual (SessionOngoing)
	err := sessionSvc.CheckOut(ctx, f.SessionOngoing.ID, parkingsession.CheckOutParams{})
	if err != nil {
		t.Fatalf("CheckOut failed: %v", err)
	}

	got, err := sessionSvc.GetByID(ctx, f.SessionOngoing.ID)
	if err != nil {
		t.Fatalf("GetByID after CheckOut failed: %v", err)
	}
	if got.Status != "completed" {
		t.Errorf("Status: got %q, want %q", got.Status, "completed")
	}
}

func TestSessionCheckOut_AlreadyCompleted(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// SessionCompleted is already checked out in fixtures
	err := sessionSvc.CheckOut(context.Background(), f.SessionCompleted.ID, parkingsession.CheckOutParams{})
	if err == nil {
		t.Error("expected error on double checkout, got nil")
	}
}

// --- GetByID ---

func TestSessionGetByID_Completed(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// SessionCompleted: CardActive, status completed
	got, err := sessionSvc.GetByID(context.Background(), f.SessionCompleted.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.CardUID != f.CardActive.CardUID {
		t.Errorf("CardUID: got %q, want %q", got.CardUID, f.CardActive.CardUID)
	}
	if got.Status != "completed" {
		t.Errorf("Status: got %q, want %q", got.Status, "completed")
	}
	if got.CheckOutTime == nil {
		t.Error("expected CheckOutTime to be set on completed session")
	}
}

func TestSessionGetByID_Ongoing(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// SessionOngoing: CardCasual, ongoing
	got, err := sessionSvc.GetByID(context.Background(), f.SessionOngoing.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Status != "ongoing" {
		t.Errorf("Status: got %q, want %q", got.Status, "ongoing")
	}
	if got.CheckOutTime != nil {
		t.Error("expected CheckOutTime to be nil on ongoing session")
	}
}

func TestSessionGetByID_NotFound(t *testing.T) {
	_, _, sessionSvc, _, cleanup := setupWithFixtures(t)
	defer cleanup()

	_, err := sessionSvc.GetByID(context.Background(), 999999)
	if err == nil {
		t.Error("expected error for non-existent session, got nil")
	}
}

// --- ListByCard ---

func TestSessionListByCard(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardActive has exactly 1 session (completed)
	sessions, err := sessionSvc.ListByCard(context.Background(), f.CardActive.CardUID)
	if err != nil {
		t.Fatalf("ListByCard failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Errorf("expected 1 session for CardActive, got %d", len(sessions))
	}
	if sessions[0].Status != "completed" {
		t.Errorf("Status: got %q, want %q", sessions[0].Status, "completed")
	}
}

func TestSessionListByCard_Multiple(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// CardActive already has 1 completed session — check in again for a second
	_, err := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: f.CardActive.CardUID,
	})
	if err != nil {
		t.Fatalf("CheckIn failed: %v", err)
	}

	sessions, err := sessionSvc.ListByCard(ctx, f.CardActive.CardUID)
	if err != nil {
		t.Fatalf("ListByCard failed: %v", err)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions for CardActive, got %d", len(sessions))
	}
}

// --- Delete ---

func TestSessionDelete(t *testing.T) {
	_, _, sessionSvc, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	err := sessionSvc.Delete(ctx, f.SessionCompleted.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = sessionSvc.GetByID(ctx, f.SessionCompleted.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestSessionDelete_NotFound(t *testing.T) {
	_, _, sessionSvc, _, cleanup := setupWithFixtures(t)
	defer cleanup()

	err := sessionSvc.Delete(context.Background(), 999999)
	if err == nil {
		t.Error("expected error for non-existent session, got nil")
	}
}

