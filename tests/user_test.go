package integration_test

import (
	"context"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/member"
	"github.com/RodrickSia/bikeKeeper/internal/user"
)

func TestUserCreate(t *testing.T) {
	_, _, _, userSvc, _, _, cleanup := setup(t)
	defer cleanup()

	u, err := userSvc.Create(context.Background(), user.CreateParams{
		Email:    "new@test.com",
		Password: "password123",
		Role:     "staff",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.ID == "" {
		t.Error("expected non-empty ID")
	}
	if u.Email != "new@test.com" {
		t.Errorf("Email: got %q, want %q", u.Email, "new@test.com")
	}
	if u.Role != "staff" {
		t.Errorf("Role: got %q, want %q", u.Role, "staff")
	}
	if u.PasswordHash == "" {
		t.Error("expected non-empty PasswordHash")
	}
}

func TestUserCreate_WithMember(t *testing.T) {
	memberSvc, _, _, userSvc, _, _, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()
	m, err := memberSvc.Create(ctx, member.CreateParams{StudentID: "99100", FullName: "Test Student"})
	if err != nil {
		t.Fatalf("create member: %v", err)
	}

	u, err := userSvc.Create(ctx, user.CreateParams{
		Email:    "student@test.com",
		Password: "password123",
		Role:     "student",
		MemberID: &m.ID,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.MemberID == nil || *u.MemberID != m.ID {
		t.Errorf("MemberID: got %v, want %q", u.MemberID, m.ID)
	}
}

func TestUserCreate_InvalidRole(t *testing.T) {
	_, _, _, userSvc, _, _, cleanup := setup(t)
	defer cleanup()

	_, err := userSvc.Create(context.Background(), user.CreateParams{
		Email:    "bad@test.com",
		Password: "password123",
		Role:     "superuser",
	})
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestUserCreate_DuplicateEmail(t *testing.T) {
	_, _, _, userSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	_, err := userSvc.Create(context.Background(), user.CreateParams{
		Email:    f.UserAdmin.Email,
		Password: "password123",
		Role:     "staff",
	})
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestUserGetByID(t *testing.T) {
	_, _, _, userSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	got, err := userSvc.GetByID(context.Background(), f.UserAdmin.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Email != "admin@test.com" {
		t.Errorf("Email: got %q, want %q", got.Email, "admin@test.com")
	}
	if got.Role != "admin" {
		t.Errorf("Role: got %q, want %q", got.Role, "admin")
	}
}

func TestUserGetByID_NotFound(t *testing.T) {
	_, _, _, userSvc, _, _, cleanup := setup(t)
	defer cleanup()

	_, err := userSvc.GetByID(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestUserGetByEmail(t *testing.T) {
	_, _, _, userSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	got, err := userSvc.GetByEmail(context.Background(), f.UserStaff.Email)
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if got.ID != f.UserStaff.ID {
		t.Errorf("ID: got %q, want %q", got.ID, f.UserStaff.ID)
	}
}

func TestUserList(t *testing.T) {
	_, _, _, userSvc, _, _, _, cleanup := setupWithFixtures(t)
	defer cleanup()

	users, err := userSvc.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUserDelete(t *testing.T) {
	_, _, _, userSvc, _, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	ctx := context.Background()
	err := userSvc.Delete(ctx, f.UserStaff.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = userSvc.GetByID(ctx, f.UserStaff.ID)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestUserDelete_NotFound(t *testing.T) {
	_, _, _, userSvc, _, _, cleanup := setup(t)
	defer cleanup()

	err := userSvc.Delete(context.Background(), "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestAuthLogin_Integration(t *testing.T) {
	_, _, _, _, authSvc, _, f, cleanup := setupWithFixtures(t)
	defer cleanup()

	token, err := authSvc.Login(context.Background(), "admin@test.com", "adminpass1")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := authSvc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != f.UserAdmin.ID {
		t.Errorf("UserID: got %q, want %q", claims.UserID, f.UserAdmin.ID)
	}
	if claims.Role != "admin" {
		t.Errorf("Role: got %q, want %q", claims.Role, "admin")
	}
}

func TestAuthLogin_WrongPassword_Integration(t *testing.T) {
	_, _, _, _, authSvc, _, _, cleanup := setupWithFixtures(t)
	defer cleanup()

	_, err := authSvc.Login(context.Background(), "admin@test.com", "wrongpassword")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}
