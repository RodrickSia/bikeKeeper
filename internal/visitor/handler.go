package visitor

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
		VisitorName  string  `json:"visitorName"`
		VisitorPhone *string `json:"visitorPhone"`
		VehiclePlate string  `json:"vehiclePlate"`
		ValidDate    string  `json:"validDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.VisitorName == "" || body.VehiclePlate == "" || body.ValidDate == "" {
		writeError(w, http.StatusBadRequest, "visitorName, vehiclePlate, and validDate are required")
		return
	}

	pass, err := h.svc.Create(r.Context(), CreateParams{
		UserID:       claims.UserID,
		VisitorName:  body.VisitorName,
		VisitorPhone: body.VisitorPhone,
		VehiclePlate: body.VehiclePlate,
		ValidDate:    body.ValidDate,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, pass)
}

func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	passes, err := h.svc.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list visitor passes")
		return
	}
	writeJSON(w, http.StatusOK, passes)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	pass, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pass)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Cancel(r.Context(), r.PathValue("id")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
