package vehicle

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	var body struct {
		LicensePlate string `json:"licensePlate"`
		Brand        string `json:"brand"`
		Model        string `json:"model"`
		Color        string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	v, err := h.svc.Add(r.Context(), CreateParams{LicensePlate: body.LicensePlate, Brand: body.Brand, Model: body.Model, Color: body.Color, OwnerID: claims.UserID})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (h *Handler) myVehicles(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	vs, err := h.svc.MyVehicles(r.Context(), claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}
	writeJSON(w, http.StatusOK, vs)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	vs, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}
	writeJSON(w, http.StatusOK, vs)
}

func (h *Handler) findByPlate(w http.ResponseWriter, r *http.Request) {
	v, err := h.svc.FindByPlate(r.Context(), r.PathValue("plate"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if v == nil {
		writeError(w, http.StatusNotFound, "vehicle not found")
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	ownerID := claims.UserID
	if err := h.svc.Remove(r.Context(), r.PathValue("id"), ownerID); err != nil {
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
