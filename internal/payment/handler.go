package payment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type CardFinder interface {
	GetByUID(ctx context.Context, cardUID string) (*CardInfo, error)
}

type CardInfo struct {
	CardUID  string
	MemberID *string
}

type Handler struct {
	svc  *Service
	cards CardFinder
}

func NewHandler(svc *Service, cards CardFinder) *Handler {
	return &Handler{svc: svc, cards: cards}
}

func (h *Handler) deposit(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	// Check ownership
	claims := auth.GetClaims(r.Context())
	if claims != nil {
		isAdminOrFaculty := claims.Role == "faculty" || claims.Role == "admin"
		if !isAdminOrFaculty {
			card, err := h.cards.GetByUID(r.Context(), cardUID)
			if err != nil {
				writeError(w, http.StatusNotFound, "card not found")
				return
			}
			if card.MemberID == nil || *card.MemberID != claims.MemberID {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
		}
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
