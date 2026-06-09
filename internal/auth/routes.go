package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string) {
	// Public — no auth needed
	mux.HandleFunc("POST "+prefix+"/auth/login", h.login)
	mux.HandleFunc("POST "+prefix+"/auth/register", h.register)
}
