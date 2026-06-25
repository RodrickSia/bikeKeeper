package incident

import (
	"encoding/json"
	"net/http"

	"github.com/RodrickSia/bikeKeeper/internal/auth"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	var staffID *string
	if s := r.URL.Query().Get("staffId"); s != "" {
		staffID = &s
	}
	incs, err := h.svc.List(r.Context(), staffID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list incidents")
		return
	}
	writeJSON(w, http.StatusOK, incs)
}

func (h *Handler) report(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	var body struct {
		VehiclePlate *string `json:"vehiclePlate"`
		Type         string  `json:"type"`
		Description  string  `json:"description"`
		Location     *string `json:"location"`
		ReporterName string  `json:"reporterName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	inc, err := h.svc.Report(r.Context(), CreateParams{
		ReportedBy: claims.UserID, ReporterName: body.ReporterName,
		VehiclePlate: body.VehiclePlate, Type: body.Type,
		Description: body.Description, Location: body.Location,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, inc)
}

func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) {
	var body struct{ Note string `json:"note"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	inc, err := h.svc.Resolve(r.Context(), r.PathValue("id"), body.Note)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func (h *Handler) escalate(w http.ResponseWriter, r *http.Request) {
	inc, err := h.svc.Escalate(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
