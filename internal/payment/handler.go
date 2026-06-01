package payment

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	var body struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	txn, err := h.svc.Deposit(r.Context(), cardUID, body.Amount)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, txn)
}

func (h *Handler) listByCard(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	txns, err := h.svc.ListByCard(r.Context(), cardUID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}
	writeJSON(w, http.StatusOK, txns)
}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	balance, err := h.svc.GetBalance(r.Context(), cardUID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"balance": balance})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
