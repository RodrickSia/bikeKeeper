package member

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string) {
	mux.HandleFunc("POST "+prefix+"/members", h.create)
	mux.HandleFunc("GET "+prefix+"/members", h.list)
	mux.HandleFunc("GET "+prefix+"/members/{id}", h.getByID)
	mux.HandleFunc("GET "+prefix+"/members/student/{studentID}", h.getByStudentID)
	mux.HandleFunc("PUT "+prefix+"/members/{id}", h.update)
	mux.HandleFunc("DELETE "+prefix+"/members/{id}", h.delete)
}
