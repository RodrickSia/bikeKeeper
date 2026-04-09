package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/card"
)

// --- Create ---

func TestCardCreate_Monthly(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// New monthly card linked to MemberA
	c, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  "NFC-NEW-MONTHLY",
		CardType: "monthly",
		MemberID: &f.MemberA.ID,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.CardUID != "NFC-NEW-MONTHLY" {
		t.Errorf("CardUID: got %q, want %q", c.CardUID, "NFC-NEW-MONTHLY")
	}
	if c.CardType != "monthly" {
		t.Errorf("CardType: got %q, want %q", c.CardType, "monthly")
	}
	if c.Status != "active" {
		t.Errorf("Status: got %q, want %q", c.Status, "active")
	}
	if c.IsInside {
		t.Error("expected IsInside=false on creation")
	}
	if c.MemberID == nil || *c.MemberID != f.MemberA.ID {
		t.Errorf("MemberID: got %v, want %q", c.MemberID, f.MemberA.ID)
	}
}

func TestCardCreate_Casual(t *testing.T) {
	_, cardSvc, _, cleanup := setup(t)
	defer cleanup()

	c, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  "NFC-CASUAL-NEW",
		CardType: "casual",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.MemberID != nil {
		t.Errorf("expected MemberID nil for casual card, got %q", *c.MemberID)
	}
}

func TestCardCreate_InvalidType(t *testing.T) {
	_, cardSvc, _, cleanup := setup(t)
	defer cleanup()

	_, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  "NFC-INVALID",
		CardType: "vip",
	})
	if err == nil {
		t.Error("expected error for invalid card type, got nil")
	}
}

// --- Read ---

func TestCardGetByUID(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardActive: NFC-MONTHLY-001, monthly, MemberA, active
	got, err := cardSvc.GetByUID(context.Background(), f.CardActive.CardUID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if got.CardUID != "NFC-MONTHLY-001" {
		t.Errorf("CardUID: got %q, want %q", got.CardUID, "NFC-MONTHLY-001")
	}
	if got.CardType != "monthly" {
		t.Errorf("CardType: got %q, want %q", got.CardType, "monthly")
	}
	if got.Status != "active" {
		t.Errorf("Status: got %q, want %q", got.Status, "active")
	}
	if got.MemberID == nil || *got.MemberID != f.MemberA.ID {
		t.Errorf("MemberID: got %v, want %q", got.MemberID, f.MemberA.ID)
	}
}

func TestCardGetByUID_Blocked(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardBlocked: NFC-MONTHLY-002, monthly, MemberB, blocked
	got, err := cardSvc.GetByUID(context.Background(), f.CardBlocked.CardUID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if got.Status != "blocked" {
		t.Errorf("Status: got %q, want %q", got.Status, "blocked")
	}
}

func TestCardGetByUID_NotFound(t *testing.T) {
	_, cardSvc, _, cleanup := setup(t)
	defer cleanup()

	_, err := cardSvc.GetByUID(context.Background(), "NONEXISTENT-CARD")
	if err == nil {
		t.Error("expected error for non-existent card, got nil")
	}
}

func TestCardListByMember(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// MemberA owns only CardActive (NFC-MONTHLY-001)
	cards, err := cardSvc.ListByMember(context.Background(), f.MemberA.ID)
	if err != nil {
		t.Fatalf("ListByMember failed: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("expected 1 card for MemberA, got %d", len(cards))
	}
	if cards[0].CardUID != f.CardActive.CardUID {
		t.Errorf("CardUID: got %q, want %q", cards[0].CardUID, f.CardActive.CardUID)
	}
}

// --- Update ---

func TestCardUpdate_Status(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// Block CardCasual (currently active)
	blocked := "blocked"
	updated, err := cardSvc.Update(context.Background(), card.UpdateParams{
		CardUID: f.CardCasual.CardUID,
		Status:  &blocked,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Status != "blocked" {
		t.Errorf("Status: got %q, want %q", updated.Status, "blocked")
	}
	if updated.CardType != "casual" {
		t.Errorf("CardType: got %q, want %q (should be unchanged)", updated.CardType, "casual")
	}
}

func TestCardUpdate_Type(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// Upgrade CardCasual to monthly
	monthly := "monthly"
	updated, err := cardSvc.Update(context.Background(), card.UpdateParams{
		CardUID:  f.CardCasual.CardUID,
		CardType: &monthly,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.CardType != "monthly" {
		t.Errorf("CardType: got %q, want %q", updated.CardType, "monthly")
	}
	if updated.Status != "active" {
		t.Errorf("Status: got %q, want %q (should be unchanged)", updated.Status, "active")
	}
}

func TestCardUpdate_InvalidStatus(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	bad := "expired"
	_, err := cardSvc.Update(context.Background(), card.UpdateParams{
		CardUID: f.CardCasual.CardUID,
		Status:  &bad,
	})
	if err == nil {
		t.Error("expected error for invalid status 'expired', got nil")
	}
}

// --- ToggleInside ---

func TestCardToggleInside(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	ctx := context.Background()

	// CardActive starts as IsInside=false
	toggled, err := cardSvc.ToggleInside(ctx, f.CardActive.CardUID)
	if err != nil {
		t.Fatalf("first ToggleInside failed: %v", err)
	}
	if !toggled.IsInside {
		t.Error("expected IsInside=true after first toggle")
	}

	toggled, err = cardSvc.ToggleInside(ctx, f.CardActive.CardUID)
	if err != nil {
		t.Fatalf("second ToggleInside failed: %v", err)
	}
	if toggled.IsInside {
		t.Error("expected IsInside=false after second toggle")
	}
}

func TestCardToggleInside_BlockedCard(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardBlocked has status "blocked" — toggle must be rejected
	_, err := cardSvc.ToggleInside(context.Background(), f.CardBlocked.CardUID)
	if err == nil {
		t.Error("expected error when toggling blocked card, got nil")
	}
}

func TestCardToggleInside_LostCard(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardLost has status "lost" — toggle must be rejected
	_, err := cardSvc.ToggleInside(context.Background(), f.CardLost.CardUID)
	if err == nil {
		t.Error("expected error when toggling lost card, got nil")
	}
}

// --- Delete ---

func TestCardDelete(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardBlocked has no sessions — safe to delete
	err := cardSvc.Delete(context.Background(), f.CardBlocked.CardUID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = cardSvc.GetByUID(context.Background(), f.CardBlocked.CardUID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestCardDelete_WithSessions(t *testing.T) {
	_, cardSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// CardActive has a completed session — ON DELETE RESTRICT prevents deletion
	err := cardSvc.Delete(context.Background(), f.CardActive.CardUID)
	if err == nil {
		t.Error("expected error when deleting card that has sessions, got nil")
	}
}

func TestCardDelete_NotFound(t *testing.T) {
	_, cardSvc, _, cleanup := setup(t)
	defer cleanup()

	err := cardSvc.Delete(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Error("expected error for non-existent card, got nil")
	}
}

