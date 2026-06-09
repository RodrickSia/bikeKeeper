package incident

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/incidents", wrap(h.list, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/incidents", wrap(h.report, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/incidents/{id}/resolve", wrap(h.resolve, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/incidents/{id}/escalate", wrap(h.escalate, authenticated, staffOnly))
}
