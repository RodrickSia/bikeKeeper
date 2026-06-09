package notification

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	ns, err := h.svc.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}
	writeJSON(w, http.StatusOK, ns)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body CreateParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	n, err := h.svc.Create(r.Context(), body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, n)
}

func (h *Handler) markRead(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.MarkRead(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) markAllRead(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if err := h.svc.MarkAllRead(r.Context(), claims.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
