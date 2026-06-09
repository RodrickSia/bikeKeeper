package support

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	var body struct {
		Category    string `json:"category"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Category == "" || body.Subject == "" || body.Description == "" {
		writeError(w, http.StatusBadRequest, "category, subject, and description are required")
		return
	}

	t, err := h.svc.Create(r.Context(), CreateParams{
		UserID:      claims.UserID,
		UserName:    "", // Can remain blank, or we could pass user email if we wanted to
		Category:    body.Category,
		Subject:     body.Subject,
		Description: body.Description,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	tickets, err := h.svc.List(r.Context(), &claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list support tickets")
		return
	}
	writeJSON(w, http.StatusOK, tickets)
}

func (h *Handler) listAll(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.svc.List(r.Context(), nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list support tickets")
		return
	}
	writeJSON(w, http.StatusOK, tickets)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	t, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}
	t, err := h.svc.UpdateStatus(r.Context(), r.PathValue("id"), body.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) addResponse(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	isAdmin := claims.Role == "admin" || claims.Role == "staff" || claims.Role == "faculty"
	t, err := h.svc.AddResponse(r.Context(), r.PathValue("id"), claims.UserID, "", body.Message, isAdmin)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
