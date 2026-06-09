package user

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/users", wrap(h.create))
	mux.Handle("GET "+prefix+"/users", wrap(h.list))
	mux.Handle("GET "+prefix+"/users/{id}", wrap(h.getByID))
	mux.Handle("PUT "+prefix+"/users/{id}/status", wrap(h.updateStatus))
	mux.Handle("DELETE "+prefix+"/users/{id}", wrap(h.delete))
}
