package monthlypass

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

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	passes, err := h.svc.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list monthly passes")
		return
	}
	if passes == nil {
		passes = []*MonthlyPass{}
	}
	writeJSON(w, http.StatusOK, passes)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		VehicleID    string  `json:"vehicleId"`
		VehiclePlate string  `json:"vehiclePlate"`
		VehicleBrand string  `json:"vehicleBrand"`
		Month        string  `json:"month"`
		StartDate    string  `json:"startDate"`
		EndDate      string  `json:"endDate"`
		Price        float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	pass, err := h.svc.Create(r.Context(), CreateParams{
		UserID:       claims.UserID,
		VehicleID:    body.VehicleID,
		VehiclePlate: body.VehiclePlate,
		VehicleBrand: body.VehicleBrand,
		Month:        body.Month,
		StartDate:    body.StartDate,
		EndDate:      body.EndDate,
		Price:        body.Price,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, pass)
}

func (h *Handler) toggleAutoRenew(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	pass, err := h.svc.ToggleAutoRenew(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pass)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
