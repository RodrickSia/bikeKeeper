package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/tests/testutil"
)

func setupMemberService(t *testing.T) (*member.Service, func()) {
	t.Helper()
	db := testutil.SetupTestDB(t)
	repo := member.NewRepository(db)
	svc := member.NewService(repo)
	cleanup := func() {
		testutil.CleanTables(t, db, "members")
	}
	cleanup()
	return svc, cleanup
}

func TestMemberCreate(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	phone := "0901234567"

	m, err := svc.Create(ctx, member.CreateParams{
		StudentID: "24520368",
		FullName:  "Nguyen Van A",
		Phone:     &phone,
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.StudentID != "24520368" {
		t.Errorf("expected StudentID '24520368', got '%s'", m.StudentID)
	}
	if m.FullName != "Nguyen Van A" {
		t.Errorf("expected FullName 'Nguyen Van A', got '%s'", m.FullName)
	}
}

func TestMemberCreate_DuplicateStudentID(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	params := member.CreateParams{
		StudentID: "24520368",
		FullName:  "Nguyen Van A",
	}

	_, err := svc.Create(ctx, params)
	if err != nil {
		t.Fatalf("first Create failed: %v", err)
	}

	_, err = svc.Create(ctx, params)
	if err == nil {
		t.Error("expected error on duplicate student_id, got nil")
	}
}

func TestMemberGetByID(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	created, _ := svc.Create(ctx, member.CreateParams{
		StudentID: "24520001",
		FullName:  "Test User",
	})

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.StudentID != "24520001" {
		t.Errorf("expected StudentID '24520001', got '%s'", got.StudentID)
	}
}

func TestMemberGetByID_NotFound(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	_, err := svc.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent ID, got nil")
	}
}

func TestMemberGetByStudentID(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	svc.Create(ctx, member.CreateParams{
		StudentID: "24520002",
		FullName:  "Student Lookup",
	})

	got, err := svc.GetByStudentID(ctx, "24520002")
	if err != nil {
		t.Fatalf("GetByStudentID failed: %v", err)
	}
	if got.FullName != "Student Lookup" {
		t.Errorf("expected FullName 'Student Lookup', got '%s'", got.FullName)
	}
}

func TestMemberList(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	svc.Create(ctx, member.CreateParams{StudentID: "10001", FullName: "A"})
	svc.Create(ctx, member.CreateParams{StudentID: "10002", FullName: "B"})
	svc.Create(ctx, member.CreateParams{StudentID: "10003", FullName: "C"})

	members, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("expected 3 members, got %d", len(members))
	}
}

func TestMemberUpdate(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	created, _ := svc.Create(ctx, member.CreateParams{
		StudentID: "24520010",
		FullName:  "Old Name",
	})

	newName := "New Name"
	newPhone := "0999999999"
	updated, err := svc.Update(ctx, member.UpdateParams{
		ID:       created.ID,
		FullName: &newName,
		Phone:    &newPhone,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.FullName != "New Name" {
		t.Errorf("expected FullName 'New Name', got '%s'", updated.FullName)
	}
	if updated.Phone == nil || *updated.Phone != "0999999999" {
		t.Error("expected phone to be updated")
	}
}

func TestMemberDelete(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	ctx := context.Background()
	created, _ := svc.Create(ctx, member.CreateParams{
		StudentID: "24520020",
		FullName:  "To Delete",
	})

	err := svc.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestMemberDelete_NotFound(t *testing.T) {
	svc, cleanup := setupMemberService(t)
	defer cleanup()

	err := svc.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent member, got nil")
	}
}
