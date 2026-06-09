package notification

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/notifications", wrap(h.listMine, authenticated))
	mux.Handle("POST "+prefix+"/notifications", wrap(h.create, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/notifications/{id}/read", wrap(h.markRead, authenticated))
	mux.Handle("POST "+prefix+"/notifications/read-all", wrap(h.markAllRead, authenticated))
}
