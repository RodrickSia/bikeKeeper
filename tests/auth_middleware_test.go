package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

func TestAuthenticate_MissingHeader(t *testing.T) {
	svc, _ := newTestService(t)
	handler := auth.Authenticate(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_InvalidFormat(t *testing.T) {
	svc, _ := newTestService(t)
	handler := auth.Authenticate(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Basic abc123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	svc, _ := newTestService(t)
	handler := auth.Authenticate(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_ValidToken(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["u@test.com"] = auth.UserRecord{
		ID:           "user-1",
		PasswordHash: hashPassword(t, "password1"),
		Role:         "staff",
	}

	token, _, _ := svc.Login(context.Background(), "u@test.com", "password1")

	var gotClaims *auth.Claims
	handler := auth.Authenticate(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotClaims = auth.GetClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	if gotClaims == nil {
		t.Fatal("expected claims in context")
	}
	if gotClaims.Role != "staff" {
		t.Errorf("role: got %q, want %q", gotClaims.Role, "staff")
	}
}

func TestRequireRole_Allowed(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["admin@test.com"] = auth.UserRecord{
		ID:           "user-1",
		PasswordHash: hashPassword(t, "password1"),
		Role:         "admin",
	}

	token, _, _ := svc.Login(context.Background(), "admin@test.com", "password1")

	handler := auth.Authenticate(svc)(
		auth.RequireRole("faculty", "admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})),
	)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequireRole_Denied(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["student@test.com"] = auth.UserRecord{
		ID:           "user-2",
		PasswordHash: hashPassword(t, "password1"),
		Role:         "student",
	}

	token, _, _ := svc.Login(context.Background(), "student@test.com", "password1")

	handler := auth.Authenticate(svc)(
		auth.RequireRole("faculty", "admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("handler should not be called")
		})),
	)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestGetRole(t *testing.T) {
	svc, mock := newTestService(t)
	mock.users["f@test.com"] = auth.UserRecord{
		ID:           "user-3",
		PasswordHash: hashPassword(t, "password1"),
		Role:         "faculty",
	}

	token, _, _ := svc.Login(context.Background(), "f@test.com", "password1")

	handler := auth.Authenticate(svc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := auth.GetRole(r.Context())
		if role != "faculty" {
			t.Errorf("GetRole: got %q, want %q", role, "faculty")
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
}

func TestGetClaims_NoContext(t *testing.T) {
	claims := auth.GetClaims(context.Background())
	if claims != nil {
		t.Errorf("expected nil claims, got %+v", claims)
	}
}

func TestGetRole_NoContext(t *testing.T) {
	role := auth.GetRole(context.Background())
	if role != "" {
		t.Errorf("expected empty role, got %q", role)
	}
}
