package device

import (
	"encoding/json"
	"net/http"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	devices, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}
	writeJSON(w, http.StatusOK, devices)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	d, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	d := &Device{}
	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Create(r.Context(), d); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, d)
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string  `json:"status"`
		Notes  *string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Status == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}
	d, err := h.svc.UpdateStatus(r.Context(), r.PathValue("id"), body.Status, body.Notes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) updateNotes(w http.ResponseWriter, r *http.Request) {
	var body struct{ Notes string `json:"notes"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	d, err := h.svc.UpdateNotes(r.Context(), r.PathValue("id"), body.Notes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) getAlerts(w http.ResponseWriter, r *http.Request) {
	var deviceID *string
	if id := r.URL.Query().Get("deviceId"); id != "" {
		deviceID = &id
	}
	alerts, err := h.svc.GetAlerts(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list alerts")
		return
	}
	writeJSON(w, http.StatusOK, alerts)
}

func (h *Handler) reportFault(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	alert, err := h.svc.ReportFault(r.Context(), r.PathValue("id"), body.Message, body.Severity)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, alert)
}

func (h *Handler) resolveAlert(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.ResolveAlert(r.Context(), r.PathValue("alertId")); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
