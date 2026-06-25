package shift

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	from := queryStr(r, "from")
	to := queryStr(r, "to")
	staffID := queryStr(r, "staffId")
	shifts, err := h.svc.List(r.Context(), from, to, staffID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list shifts")
		return
	}
	writeJSON(w, http.StatusOK, shifts)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body CreateParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sh, err := h.svc.Create(r.Context(), body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, sh)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	sh, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	var body struct{ Status string `json:"status"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}
	sh, err := h.svc.UpdateStatus(r.Context(), r.PathValue("id"), body.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func (h *Handler) updateNotes(w http.ResponseWriter, r *http.Request) {
	var body struct{ Notes *string `json:"notes"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sh, err := h.svc.UpdateNotes(r.Context(), r.PathValue("id"), body.Notes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) assignStaff(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	var body struct{ StaffID string `json:"staffId"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.StaffID == "" {
		writeError(w, http.StatusBadRequest, "staffId is required")
		return
	}
	assignedBy := claims.UserID
	sh, err := h.svc.AssignStaff(r.Context(), r.PathValue("id"), body.StaffID, &assignedBy)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func (h *Handler) removeStaff(w http.ResponseWriter, r *http.Request) {
	var body struct{ StaffID string `json:"staffId"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.StaffID == "" {
		writeError(w, http.StatusBadRequest, "staffId is required")
		return
	}
	sh, err := h.svc.RemoveStaff(r.Context(), r.PathValue("id"), body.StaffID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sh)
}

func queryStr(r *http.Request, key string) *string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	return &v
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
