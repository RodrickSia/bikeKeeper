package parkingsession

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /sessions/checkin
func (h *Handler) checkIn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	cardUID := r.FormValue("cardUid")
	if cardUID == "" {
		writeError(w, http.StatusBadRequest, "cardUid is required")
		return
	}

	imgPlateIn, err := readFormFile(r, "imgPlateIn")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read imgPlateIn")
		return
	}

	imgPersonIn, err := readFormFile(r, "imgPersonIn")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read imgPersonIn")
		return
	}

	session, err := h.svc.CheckIn(r.Context(), CheckInParams{
		CardUID:     cardUID,
		ImgPlateIn:  imgPlateIn,
		ImgPersonIn: imgPersonIn,
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	imgPlateOut, err := readFormFile(r, "imgPlateOut")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read imgPlateOut")
		return
	}

	imgPersonOut, err := readFormFile(r, "imgPersonOut")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read imgPersonOut")
		return
	}

	if err := h.svc.CheckOut(r.Context(), id, CheckOutParams{
		ImgPlateOut:  imgPlateOut,
		ImgPersonOut: imgPersonOut,
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

func readFormFile(r *http.Request, field string) ([]byte, error) {
	file, _, err := r.FormFile(field)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}


