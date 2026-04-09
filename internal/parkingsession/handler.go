package parkingsession

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /sessions/checkin
func (h *Handler) checkIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CardUID         string  `json:"cardUid"`
		PlateIn         *string `json:"plateIn"`
		ImgPlateInPath  *string `json:"imgPlateInPath"`
		ImgPersonInPath *string `json:"imgPersonInPath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUid is required")
		return
	}

	session, err := h.svc.CheckIn(r.Context(), CheckInParams{
		CardUID:         body.CardUID,
		PlateIn:         body.PlateIn,
		ImgPlateInPath:  body.ImgPlateInPath,
		ImgPersonInPath: body.ImgPersonInPath,
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

// POST /sessions/{id}/checkout
func (h *Handler) checkOut(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var body struct {
		PlateOut         *string `json:"plateOut"`
		ImgPlateOutPath  *string `json:"imgPlateOutPath"`
		ImgPersonOutPath *string `json:"imgPersonOutPath"`
		Cost             string  `json:"cost"`
		IsWarning        bool    `json:"isWarning"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cost, err := decimal.NewFromString(body.Cost)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid cost value")
		return
	}

	if err := h.svc.CheckOut(r.Context(), id, CheckOutParams{
		PlateOut:         body.PlateOut,
		ImgPlateOutPath:  body.ImgPlateOutPath,
		ImgPersonOutPath: body.ImgPersonOutPath,
		Cost:             cost,
		IsWarning:        body.IsWarning,
	}); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /sessions/{id}
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	session, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, session)
}

// GET /sessions/card/{cardUID}
func (h *Handler) listByCard(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	sessions, err := h.svc.ListByCard(r.Context(), cardUID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list sessions")
		return
	}
	writeJSON(w, http.StatusOK, sessions)
}

// DELETE /sessions/{id}
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid session id")
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}


