package cardrequest

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

// POST /card-requests  (student submits)
func (h *Handler) submit(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil || claims.MemberID == "" {
		writeError(w, http.StatusForbidden, "only members can submit card requests")
		return
	}
	var body struct {
		VehiclePlate string  `json:"vehiclePlate"`
		VehicleBrand string  `json:"vehicleBrand"`
		VehicleModel string  `json:"vehicleModel"`
		VehicleColor string  `json:"vehicleColor"`
		IDCardNumber string  `json:"idCardNumber"`
		Note         *string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.VehiclePlate == "" || body.VehicleBrand == "" || body.IDCardNumber == "" {
		writeError(w, http.StatusBadRequest, "vehiclePlate, vehicleBrand, and idCardNumber are required")
		return
	}
	req, err := h.svc.Submit(r.Context(), CreateParams{
		MemberID:     claims.MemberID,
		VehiclePlate: body.VehiclePlate,
		VehicleBrand: body.VehicleBrand,
		VehicleModel: body.VehicleModel,
		VehicleColor: body.VehicleColor,
		IDCardNumber: body.IDCardNumber,
		Note:         body.Note,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, req)
}

// GET /card-requests/me  (student views own requests)
func (h *Handler) listMine(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil || claims.MemberID == "" {
		writeError(w, http.StatusForbidden, "only members can view their requests")
		return
	}
	reqs, err := h.svc.ListByMember(r.Context(), claims.MemberID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list requests")
		return
	}
	writeJSON(w, http.StatusOK, reqs)
}

// GET /card-requests  (admin/faculty lists all, optional ?status=pending)
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	var status *string
	if s := r.URL.Query().Get("status"); s != "" {
		status = &s
	}
	reqs, err := h.svc.List(r.Context(), status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list requests")
		return
	}
	writeJSON(w, http.StatusOK, reqs)
}

// POST /card-requests/{id}/approve
func (h *Handler) approve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	claims := auth.GetClaims(r.Context())
	var body struct {
		CardUID string `json:"cardUid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUid is required")
		return
	}
	req, err := h.svc.Approve(r.Context(), id, body.CardUID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

// POST /card-requests/{id}/reject
func (h *Handler) reject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	claims := auth.GetClaims(r.Context())
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Reason == "" {
		writeError(w, http.StatusBadRequest, "reason is required")
		return
	}
	req, err := h.svc.Reject(r.Context(), id, body.Reason, claims.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

// POST /card-requests/{id}/block
func (h *Handler) block(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, err := h.svc.Block(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
