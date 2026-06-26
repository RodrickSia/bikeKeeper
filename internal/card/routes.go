package card

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/cards", wrap(h.create))
	mux.Handle("GET "+prefix+"/cards/{cardUID}", wrap(h.getByUID))
	mux.Handle("GET "+prefix+"/cards/casual-available", wrap(h.getAvailableCasual))
	mux.Handle("GET "+prefix+"/members-cards/{memberID}", wrap(h.listByMember))
	mux.Handle("PUT "+prefix+"/cards/{cardUID}", wrap(h.update))
	mux.Handle("POST "+prefix+"/cards/{cardUID}/toggle", wrap(h.toggleInside))
	mux.Handle("DELETE "+prefix+"/cards/{cardUID}", wrap(h.delete))
}
