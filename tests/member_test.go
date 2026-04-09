package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/member"
)

// --- Create ---

func TestMemberCreate(t *testing.T) {
	memberSvc, _, _, cleanup := setup(t)
	defer cleanup()

	phone := "0908888888"
	m, err := memberSvc.Create(context.Background(), member.CreateParams{
		StudentID: "99000",
		FullName:  "New Student",
		Phone:     &phone,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.StudentID != "99000" {
		t.Errorf("StudentID: got %q, want %q", m.StudentID, "99000")
	}
	if m.FullName != "New Student" {
		t.Errorf("FullName: got %q, want %q", m.FullName, "New Student")
	}
	if m.Phone == nil || *m.Phone != "0908888888" {
		t.Error("expected Phone '0908888888'")
	}
}

func TestMemberCreate_NoPhone(t *testing.T) {
	memberSvc, _, _, cleanup := setup(t)
	defer cleanup()

	m, err := memberSvc.Create(context.Background(), member.CreateParams{
		StudentID: "99001",
		FullName:  "No Phone Student",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if m.Phone != nil {
		t.Errorf("expected Phone nil, got %q", *m.Phone)
	}
}

func TestMemberCreate_DuplicateStudentID(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// MemberA (StudentID "24520001") already exists in the DB
	_, err := memberSvc.Create(context.Background(), member.CreateParams{
		StudentID: f.MemberA.StudentID,
		FullName:  "Duplicate",
	})
	if err == nil {
		t.Error("expected error for duplicate student_id, got nil")
	}
}

// --- Read ---

func TestMemberGetByID(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	got, err := memberSvc.GetByID(context.Background(), f.MemberA.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.StudentID != "24520001" {
		t.Errorf("StudentID: got %q, want %q", got.StudentID, "24520001")
	}
	if got.FullName != "Nguyen Van A" {
		t.Errorf("FullName: got %q, want %q", got.FullName, "Nguyen Van A")
	}
	if got.Phone == nil || *got.Phone != "0901234567" {
		t.Error("expected Phone '0901234567'")
	}
}

func TestMemberGetByID_NotFound(t *testing.T) {
	memberSvc, _, _, cleanup := setup(t)
	defer cleanup()

	_, err := memberSvc.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}

func TestMemberGetByStudentID(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// MemberB has StudentID "24520002" and no phone
	got, err := memberSvc.GetByStudentID(context.Background(), f.MemberB.StudentID)
	if err != nil {
		t.Fatalf("GetByStudentID failed: %v", err)
	}
	if got.ID != f.MemberB.ID {
		t.Errorf("ID: got %q, want %q", got.ID, f.MemberB.ID)
	}
	if got.FullName != "Tran Thi B" {
		t.Errorf("FullName: got %q, want %q", got.FullName, "Tran Thi B")
	}
	if got.Phone != nil {
		t.Errorf("expected Phone nil for MemberB, got %q", *got.Phone)
	}
}

func TestMemberGetByStudentID_NotFound(t *testing.T) {
	memberSvc, _, _, cleanup := setup(t)
	defer cleanup()

	_, err := memberSvc.GetByStudentID(context.Background(), "00000000")
	if err == nil {
		t.Error("expected error for non-existent student_id, got nil")
	}
}

func TestMemberList(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	members, err := memberSvc.List(context.Background())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("expected 3 members, got %d", len(members))
	}

	ids := make(map[string]bool, 3)
	for _, m := range members {
		ids[m.ID] = true
	}
	for _, want := range []*member.Member{f.MemberA, f.MemberB, f.MemberC} {
		if !ids[want.ID] {
			t.Errorf("member %q (%s) not found in List result", want.FullName, want.ID)
		}
	}
}

// --- Update ---

func TestMemberUpdate_FullName(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// Update MemberC's full name; phone should remain unchanged
	newName := "Le Van C (Updated)"
	updated, err := memberSvc.Update(context.Background(), member.UpdateParams{
		ID:       f.MemberC.ID,
		FullName: &newName,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.FullName != newName {
		t.Errorf("FullName: got %q, want %q", updated.FullName, newName)
	}
	if updated.Phone == nil || *updated.Phone != "0972345678" {
		t.Error("expected Phone '0972345678' to remain unchanged")
	}
}

func TestMemberUpdate_AddPhone(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// MemberB initially has no phone — add one
	newPhone := "0911222333"
	updated, err := memberSvc.Update(context.Background(), member.UpdateParams{
		ID:    f.MemberB.ID,
		Phone: &newPhone,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Phone == nil || *updated.Phone != newPhone {
		t.Errorf("Phone: got %v, want %q", updated.Phone, newPhone)
	}
	if updated.FullName != "Tran Thi B" {
		t.Errorf("FullName: got %q, want %q (should be unchanged)", updated.FullName, "Tran Thi B")
	}
}

// --- Delete ---

func TestMemberDelete(t *testing.T) {
	memberSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	// Deleting MemberA sets CardActive.member_id → NULL (ON DELETE SET NULL)
	err := memberSvc.Delete(context.Background(), f.MemberA.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = memberSvc.GetByID(context.Background(), f.MemberA.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestMemberDelete_NotFound(t *testing.T) {
	memberSvc, _, _, cleanup := setup(t)
	defer cleanup()

	err := memberSvc.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent member, got nil")
	}
}
