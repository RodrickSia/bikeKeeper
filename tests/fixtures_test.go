package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/card"
	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/internal/parkingsession"
	"github.com/RodrickSia/bikeKeeper/tests/testutil"
	"github.com/shopspring/decimal"
)

// Fixtures holds all pre-seeded dummy data available to each test.
//
// Dummy data overview:
//
//	Members
//	  MemberA — "Nguyen Van A", student 24520001, phone 0901234567, owns CardActive
//	  MemberB — "Tran Thi B",   student 24520002, no phone,         owns CardBlocked
//	  MemberC — "Le Van C",     student 24520003, phone 0972345678, owns CardLost
//
//	Cards
//	  CardActive  — NFC-MONTHLY-001, monthly, MemberA, active  (1 completed session)
//	  CardBlocked — NFC-MONTHLY-002, monthly, MemberB, blocked (no sessions)
//	  CardLost    — NFC-MONTHLY-003, monthly, MemberC, lost    (no sessions)
//	  CardCasual  — NFC-CASUAL-001,  casual,  no member, active (1 ongoing session)
//
//	Sessions
//	  SessionCompleted — CardActive,  plate 59F1-12345, cost 5000, status completed
//	  SessionOngoing   — CardCasual,  plate 51H-99999,  status ongoing
type Fixtures struct {
	MemberA *member.Member
	MemberB *member.Member
	MemberC *member.Member

	CardActive  *card.Card
	CardBlocked *card.Card
	CardLost    *card.Card
	CardCasual  *card.Card

	SessionCompleted *parkingsession.ParkingSession
	SessionOngoing   *parkingsession.ParkingSession
}

// setup returns all three services wired to a clean, empty test DB.
func setup(t *testing.T) (*member.Service, *card.Service, *parkingsession.Service, func()) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	memberSvc, cardSvc, sessionSvc := buildServices(db)
	cleanup := func() {
		testutil.CleanTables(t, db, "parking_sessions", "cards", "members")
	}
	cleanup()
	return memberSvc, cardSvc, sessionSvc, cleanup
}

// setupWithFixtures seeds the full dummy dataset and returns services alongside
// the Fixtures struct so tests can reference known values.
func setupWithFixtures(t *testing.T) (*member.Service, *card.Service, *parkingsession.Service, *Fixtures, func()) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	memberSvc, cardSvc, sessionSvc := buildServices(db)
	cleanup := func() {
		testutil.CleanTables(t, db, "parking_sessions", "cards", "members")
	}
	cleanup()

	f, err := seedFixtures(context.Background(), memberSvc, cardSvc, sessionSvc)
	if err != nil {
		t.Fatalf("seedFixtures: %v", err)
	}
	return memberSvc, cardSvc, sessionSvc, f, cleanup
}

func buildServices(db *sql.DB) (*member.Service, *card.Service, *parkingsession.Service) {
	memberSvc := member.NewService(member.NewRepository(db))
	cardSvc := card.NewService(card.NewRepository(db))
	sessionSvc := parkingsession.NewService(parkingsession.NewRepository(db))
	return memberSvc, cardSvc, sessionSvc
}

func seedFixtures(
	ctx context.Context,
	memberSvc *member.Service,
	cardSvc *card.Service,
	sessionSvc *parkingsession.Service,
) (*Fixtures, error) {
	// --- Members ---
	phone1 := "0901234567"
	memberA, err := memberSvc.Create(ctx, member.CreateParams{
		StudentID: "24520001",
		FullName:  "Nguyen Van A",
		Phone:     &phone1,
	})
	if err != nil {
		return nil, fmt.Errorf("create memberA: %w", err)
	}

	memberB, err := memberSvc.Create(ctx, member.CreateParams{
		StudentID: "24520002",
		FullName:  "Tran Thi B",
	})
	if err != nil {
		return nil, fmt.Errorf("create memberB: %w", err)
	}

	phone3 := "0972345678"
	memberC, err := memberSvc.Create(ctx, member.CreateParams{
		StudentID: "24520003",
		FullName:  "Le Van C",
		Phone:     &phone3,
	})
	if err != nil {
		return nil, fmt.Errorf("create memberC: %w", err)
	}

	// --- Cards ---
	cardActive, err := cardSvc.Create(ctx, card.CreateParams{
		CardUID:  "NFC-MONTHLY-001",
		CardType: "monthly",
		MemberID: &memberA.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create cardActive: %w", err)
	}

	rawBlocked, err := cardSvc.Create(ctx, card.CreateParams{
		CardUID:  "NFC-MONTHLY-002",
		CardType: "monthly",
		MemberID: &memberB.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create cardBlocked: %w", err)
	}
	blockedStatus := "blocked"
	cardBlocked, err := cardSvc.Update(ctx, card.UpdateParams{
		CardUID: rawBlocked.CardUID,
		Status:  &blockedStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("block cardBlocked: %w", err)
	}

	rawLost, err := cardSvc.Create(ctx, card.CreateParams{
		CardUID:  "NFC-MONTHLY-003",
		CardType: "monthly",
		MemberID: &memberC.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("create cardLost: %w", err)
	}
	lostStatus := "lost"
	cardLost, err := cardSvc.Update(ctx, card.UpdateParams{
		CardUID: rawLost.CardUID,
		Status:  &lostStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("mark cardLost: %w", err)
	}

	cardCasual, err := cardSvc.Create(ctx, card.CreateParams{
		CardUID:  "NFC-CASUAL-001",
		CardType: "casual",
	})
	if err != nil {
		return nil, fmt.Errorf("create cardCasual: %w", err)
	}

	// --- Sessions ---
	plateIn1 := "59F1-12345"
	sessionCompleted, err := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: cardActive.CardUID,
		PlateIn: &plateIn1,
	})
	if err != nil {
		return nil, fmt.Errorf("check-in sessionCompleted: %w", err)
	}
	plateOut1 := "59F1-12345"
	if err := sessionSvc.CheckOut(ctx, sessionCompleted.ID, parkingsession.CheckOutParams{
		PlateOut: &plateOut1,
		Cost:     decimal.NewFromInt(5000),
	}); err != nil {
		return nil, fmt.Errorf("check-out sessionCompleted: %w", err)
	}
	sessionCompleted, err = sessionSvc.GetByID(ctx, sessionCompleted.ID)
	if err != nil {
		return nil, fmt.Errorf("re-fetch sessionCompleted: %w", err)
	}

	plateIn2 := "51H-99999"
	sessionOngoing, err := sessionSvc.CheckIn(ctx, parkingsession.CheckInParams{
		CardUID: cardCasual.CardUID,
		PlateIn: &plateIn2,
	})
	if err != nil {
		return nil, fmt.Errorf("check-in sessionOngoing: %w", err)
	}

	return &Fixtures{
		MemberA:          memberA,
		MemberB:          memberB,
		MemberC:          memberC,
		CardActive:       cardActive,
		CardBlocked:      cardBlocked,
		CardLost:         cardLost,
		CardCasual:       cardCasual,
		SessionCompleted: sessionCompleted,
		SessionOngoing:   sessionOngoing,
	}, nil
}
