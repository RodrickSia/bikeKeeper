package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type UserCreator interface {
	CreateRegister(ctx context.Context, email, password, role string, memberID *string) error
}

type MemberCreator interface {
	CreateRegister(ctx context.Context, studentID, fullName string, phone *string) (string, error)
}

type Handler struct {
	svc           *Service
	userCreator   UserCreator
	memberCreator MemberCreator
}

func NewHandler(svc *Service, uc UserCreator, mc MemberCreator) *Handler {
	return &Handler{
		svc:           svc,
		userCreator:   uc,
		memberCreator: mc,
	}
}

// POST /auth/login
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Email == "" || body.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	token, userResp, err := h.svc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		log.Printf("Login failed for email %s: %v", body.Email, err)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  userResp,
	})
}

// POST /auth/register
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string  `json:"name"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Role      string  `json:"role"`
		StudentID string  `json:"studentId"`
		Phone     *string `json:"phone,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Email == "" || body.Password == "" || body.Role == "" || body.Name == "" {
		writeError(w, http.StatusBadRequest, "name, email, password, and role are required")
		return
	}
	if len(body.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	var memberID *string
	// If role is student or if StudentID is provided, create the Member profile
	if body.Role == "student" || body.StudentID != "" {
		if body.StudentID == "" {
			writeError(w, http.StatusBadRequest, "studentId is required for student role")
			return
		}
		mID, err := h.memberCreator.CreateRegister(r.Context(), body.StudentID, body.Name, body.Phone)
		if err != nil {
			writeError(w, http.StatusConflict, "failed to create member profile: "+err.Error())
			return
		}
		memberID = &mID
	}

	// Create user with default pending status
	err := h.userCreator.CreateRegister(r.Context(), body.Email, body.Password, body.Role, memberID)
	if err != nil {
		writeError(w, http.StatusConflict, "failed to register user: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Registration successful, pending admin approval"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
