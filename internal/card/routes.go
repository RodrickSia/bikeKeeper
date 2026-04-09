package card

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string) {
	mux.HandleFunc("POST "+prefix+"/cards", h.create)
	mux.HandleFunc("GET "+prefix+"/cards/{cardUID}", h.getByUID)
	mux.HandleFunc("GET "+prefix+"/cards/member/{memberID}", h.listByMember)
	mux.HandleFunc("PUT "+prefix+"/cards/{cardUID}", h.update)
	mux.HandleFunc("POST "+prefix+"/cards/{cardUID}/toggle", h.toggleInside)
	mux.HandleFunc("DELETE "+prefix+"/cards/{cardUID}", h.delete)
}
