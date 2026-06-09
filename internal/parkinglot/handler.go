package parkinglot

import (
	"encoding/json"
	"net/http"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	lots, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list parking lots")
		return
	}
	writeJSON(w, http.StatusOK, lots)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	lot, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, lot)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name          string  `json:"name"`
		Address       string  `json:"address"`
		Type          string  `json:"type"`
		TotalCapacity int     `json:"totalCapacity"`
		OpenTime      string  `json:"openTime"`
		CloseTime     string  `json:"closeTime"`
		ContactPhone  *string `json:"contactPhone"`
		ManagerName   *string `json:"managerName"`
		Description   *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	lot, err := h.svc.Create(r.Context(), CreateParams(body))
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, lot)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var body UpdateParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	lot, err := h.svc.Update(r.Context(), r.PathValue("id"), body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, lot)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
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
