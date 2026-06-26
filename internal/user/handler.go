package user

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /users (create user — faculty + admin only)
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		Role     string  `json:"role"`
		MemberID *string `json:"memberId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Email == "" || body.Password == "" || body.Role == "" {
		writeError(w, http.StatusBadRequest, "email, password, and role are required")
		return
	}
	if len(body.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	callerRole := auth.GetRole(r.Context())
	if callerRole != RoleAdmin && callerRole != RoleFaculty {
		writeError(w, http.StatusForbidden, "only faculty and admin can create users")
		return
	}
	if body.Role == RoleAdmin && callerRole != RoleAdmin {
		writeError(w, http.StatusForbidden, "only admins can create admin users")
		return
	}

	user, err := h.svc.Create(r.Context(), CreateParams{
		Email:    body.Email,
		Password: body.Password,
		Role:     body.Role,
		MemberID: body.MemberID,
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// GET /users
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GET /users/{id}
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// DELETE /users/{id}
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	callerRole := auth.GetRole(r.Context())
	if callerRole != RoleAdmin && callerRole != RoleFaculty {
		writeError(w, http.StatusForbidden, "only faculty and admin can delete users")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PUT /users/{id}/status
func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	callerRole := auth.GetRole(r.Context())
	if callerRole != RoleAdmin && callerRole != RoleFaculty {
		writeError(w, http.StatusForbidden, "only faculty and admin can update user status")
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, body.Status); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
