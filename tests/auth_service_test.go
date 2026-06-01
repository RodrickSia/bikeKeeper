package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type mockUserFinder struct {
	users map[string]auth.UserRecord
}

func (m *mockUserFinder) GetByEmail(_ context.Context, email string) (auth.UserRecord, error) {
	u, ok := m.users[email]
	if !ok {
		return auth.UserRecord{}, fmt.Errorf("not found")
	}
	return u, nil
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	return string(h)
}

func newTestService(t *testing.T) (*auth.Service, *mockUserFinder) {
	t.Helper()
	mock := &mockUserFinder{users: map[string]auth.UserRecord{}}
	t.Setenv("JWT_SECRET", "test-secret")
	svc := auth.NewService(mock)
	return svc, mock
}

func TestLogin_Success(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["staff@test.com"] = auth.UserRecord{
		ID:           "user-1",
		PasswordHash: hashPassword(t, "password123"),
		Role:         "staff",
	}

	token, err := svc.Login(context.Background(), "staff@test.com", "password123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["staff@test.com"] = auth.UserRecord{
		ID:           "user-1",
		PasswordHash: hashPassword(t, "password123"),
		Role:         "staff",
	}

	_, err := svc.Login(context.Background(), "staff@test.com", "wrong")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.Login(context.Background(), "nobody@test.com", "password123")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	svc, mock := newTestService(t)
	memberID := "member-1"
	mock.users["student@test.com"] = auth.UserRecord{
		ID:           "user-2",
		PasswordHash: hashPassword(t, "pass1234"),
		Role:         "student",
		MemberID:     &memberID,
	}

	token, err := svc.Login(context.Background(), "student@test.com", "pass1234")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != "user-2" {
		t.Errorf("UserID: got %q, want %q", claims.UserID, "user-2")
	}
	if claims.Role != "student" {
		t.Errorf("Role: got %q, want %q", claims.Role, "student")
	}
	if claims.MemberID != "member-1" {
		t.Errorf("MemberID: got %q, want %q", claims.MemberID, "member-1")
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.ValidateToken("garbage.token.value")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateToken_TamperedSecret(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["a@test.com"] = auth.UserRecord{
		ID:           "user-3",
		PasswordHash: hashPassword(t, "pass1234"),
		Role:         "admin",
	}

	token, err := svc.Login(context.Background(), "a@test.com", "pass1234")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	t.Setenv("JWT_SECRET", "different-secret")
	svc2 := auth.NewService(mock)

	_, err = svc2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for token signed with different secret")
	}
}
