package parkingsession

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string) {
	mux.HandleFunc("POST "+prefix+"/sessions/checkin", h.checkIn)
	mux.HandleFunc("POST "+prefix+"/sessions/{id}/checkout", h.checkOut)
	mux.HandleFunc("GET "+prefix+"/sessions/{id}", h.getByID)
	mux.HandleFunc("GET "+prefix+"/sessions/card/{cardUID}", h.listByCard)
	mux.HandleFunc("DELETE "+prefix+"/sessions/{id}", h.delete)
}
