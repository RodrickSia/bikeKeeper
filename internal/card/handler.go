package card

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

// POST /cards
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CardUID  string  `json:"cardUid"`
		CardType string  `json:"cardType"`
		MemberID *string `json:"memberId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.CardUID == "" || body.CardType == "" {
		writeError(w, http.StatusBadRequest, "cardUid and cardType are required")
		return
	}

	card, err := h.svc.Create(r.Context(), CreateParams{
		CardUID:  body.CardUID,
		CardType: body.CardType,
		MemberID: body.MemberID,
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, card)
}

// GET /cards/{cardUID}
func (h *Handler) getByUID(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	card, err := h.svc.GetByUID(r.Context(), cardUID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

// GET /cards/member/{memberID}
func (h *Handler) listByMember(w http.ResponseWriter, r *http.Request) {
	memberID := r.PathValue("memberID")
	if memberID == "" {
		writeError(w, http.StatusBadRequest, "memberID is required")
		return
	}

	cards, err := h.svc.ListByMember(r.Context(), memberID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list cards")
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

// PUT /cards/{cardUID}
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	var body struct {
		CardType *string `json:"cardType"`
		MemberID *string `json:"memberId"`
		Status   *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	card, err := h.svc.Update(r.Context(), UpdateParams{
		CardUID:  cardUID,
		CardType: body.CardType,
		MemberID: body.MemberID,
		Status:   body.Status,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

// POST /cards/{cardUID}/toggle
func (h *Handler) toggleInside(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	card, err := h.svc.ToggleInside(r.Context(), cardUID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

// DELETE /cards/{cardUID}
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	cardUID := r.PathValue("cardUID")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUID is required")
		return
	}

	if err := h.svc.Delete(r.Context(), cardUID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
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
