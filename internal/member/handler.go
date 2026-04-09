package member

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

// POST /members
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		StudentID string  `json:"studentId"`
		FullName  string  `json:"fullName"`
		Phone     *string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.StudentID == "" || body.FullName == "" {
		writeError(w, http.StatusBadRequest, "studentId and fullName are required")
		return
	}

	member, err := h.svc.Create(r.Context(), CreateParams{
		StudentID: body.StudentID,
		FullName:  body.FullName,
		Phone:     body.Phone,
	})
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

// GET /members/{id}
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	member, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, member)
}

// GET /members/student/{studentID}
func (h *Handler) getByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID := r.PathValue("studentID")
	if studentID == "" {
		writeError(w, http.StatusBadRequest, "studentID is required")
		return
	}

	member, err := h.svc.GetByStudentID(r.Context(), studentID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, member)
}

// GET /members
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	members, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list members")
		return
	}
	writeJSON(w, http.StatusOK, members)
}

// PUT /members/{id}
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var body struct {
		FullName *string `json:"fullName"`
		Phone    *string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	member, err := h.svc.Update(r.Context(), UpdateParams{
		ID:       id,
		FullName: body.FullName,
		Phone:    body.Phone,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, member)
}

// DELETE /members/{id}
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
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
