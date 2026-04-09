package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/card"
	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/tests/testutil"
)

func setupCardService(t *testing.T) (*card.Service, *member.Service, func()) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	cardRepo := card.NewRepository(db)
	cardSvc := card.NewService(cardRepo)
	memberRepo := member.NewRepository(db)
	memberSvc := member.NewService(memberRepo)
	cleanup := func() {
		testutil.CleanTables(t, db, "cards", "members")
	}
	cleanup()
	return cardSvc, memberSvc, cleanup
}

func createTestMember(t *testing.T, memberSvc *member.Service, studentID string) *member.Member {
	t.Helper()
	m, err := memberSvc.Create(context.Background(), member.CreateParams{
		StudentID: studentID,
		FullName:  "Test Member " + studentID,
	})
	if err != nil {
		t.Fatalf("failed to create test member: %v", err)
	}
	return m
}

func TestCardCreate(t *testing.T) {
	cardSvc, memberSvc, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	m := createTestMember(t, memberSvc, "99001")

	c, err := cardSvc.Create(ctx, card.CreateParams{
		CardUID:  "NFC-001",
		CardType: "monthly",
		MemberID: &m.ID,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if c.CardUID != "NFC-001" {
		t.Errorf("expected CardUID 'NFC-001', got '%s'", c.CardUID)
	}
	if c.CardType != "monthly" {
		t.Errorf("expected CardType 'monthly', got '%s'", c.CardType)
	}
	if c.Status != "active" {
		t.Errorf("expected Status 'active', got '%s'", c.Status)
	}
	if c.IsInside {
		t.Error("expected IsInside false on creation")
	}
}

func TestCardCreate_Casual(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	c, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  "NFC-CASUAL-001",
		CardType: "casual",
	})
	if err != nil {
		t.Fatalf("Create casual failed: %v", err)
	}
	if c.MemberID != nil {
		t.Error("expected nil MemberID for casual card")
	}
}

func TestCardCreate_InvalidType(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	_, err := cardSvc.Create(context.Background(), card.CreateParams{
		CardUID:  "NFC-BAD",
		CardType: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid card type, got nil")
	}
}

func TestCardGetByUID(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-GET-001", CardType: "casual"})

	got, err := cardSvc.GetByUID(ctx, "NFC-GET-001")
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if got.CardUID != "NFC-GET-001" {
		t.Errorf("expected 'NFC-GET-001', got '%s'", got.CardUID)
	}
}

func TestCardGetByUID_NotFound(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	_, err := cardSvc.GetByUID(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Error("expected error for non-existent card, got nil")
	}
}

func TestCardListByMember(t *testing.T) {
	cardSvc, memberSvc, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	m := createTestMember(t, memberSvc, "99002")

	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-M1", CardType: "monthly", MemberID: &m.ID})
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-M2", CardType: "monthly", MemberID: &m.ID})

	cards, err := cardSvc.ListByMember(ctx, m.ID)
	if err != nil {
		t.Fatalf("ListByMember failed: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(cards))
	}
}

func TestCardUpdate(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-UPD", CardType: "casual"})

	newType := "monthly"
	newStatus := "blocked"
	updated, err := cardSvc.Update(ctx, card.UpdateParams{
		CardUID:  "NFC-UPD",
		CardType: &newType,
		Status:   &newStatus,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.CardType != "monthly" {
		t.Errorf("expected CardType 'monthly', got '%s'", updated.CardType)
	}
	if updated.Status != "blocked" {
		t.Errorf("expected Status 'blocked', got '%s'", updated.Status)
	}
}

func TestCardUpdate_InvalidStatus(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-BAD-STATUS", CardType: "casual"})

	badStatus := "expired"
	_, err := cardSvc.Update(ctx, card.UpdateParams{
		CardUID: "NFC-BAD-STATUS",
		Status:  &badStatus,
	})
	if err == nil {
		t.Error("expected error for invalid status, got nil")
	}
}

func TestCardToggleInside(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-TOGGLE", CardType: "casual"})

	// toggle in
	toggled, err := cardSvc.ToggleInside(ctx, "NFC-TOGGLE")
	if err != nil {
		t.Fatalf("first ToggleInside failed: %v", err)
	}
	if !toggled.IsInside {
		t.Error("expected IsInside true after first toggle")
	}

	// toggle out
	toggled, err = cardSvc.ToggleInside(ctx, "NFC-TOGGLE")
	if err != nil {
		t.Fatalf("second ToggleInside failed: %v", err)
	}
	if toggled.IsInside {
		t.Error("expected IsInside false after second toggle")
	}
}

func TestCardToggleInside_BlockedCard(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-BLOCKED", CardType: "casual"})

	blocked := "blocked"
	cardSvc.Update(ctx, card.UpdateParams{CardUID: "NFC-BLOCKED", Status: &blocked})

	_, err := cardSvc.ToggleInside(ctx, "NFC-BLOCKED")
	if err == nil {
		t.Error("expected error when toggling blocked card, got nil")
	}
}

func TestCardDelete(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	ctx := context.Background()
	cardSvc.Create(ctx, card.CreateParams{CardUID: "NFC-DEL", CardType: "casual"})

	err := cardSvc.Delete(ctx, "NFC-DEL")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = cardSvc.GetByUID(ctx, "NFC-DEL")
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestCardDelete_NotFound(t *testing.T) {
	cardSvc, _, cleanup := setupCardService(t)
	defer cleanup()

	err := cardSvc.Delete(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Error("expected error for non-existent card, got nil")
	}
}
