package visitor

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/visitor/passes", wrap(h.create, authenticated))
	mux.Handle("GET "+prefix+"/visitor/passes/me", wrap(h.listMine, authenticated))
	mux.Handle("GET "+prefix+"/visitor/passes/{id}", wrap(h.getByID, authenticated))
	mux.Handle("POST "+prefix+"/visitor/passes/{id}/cancel", wrap(h.cancel, authenticated))
}
