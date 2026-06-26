package monthlypass

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/monthly-passes", wrap(h.list))
	mux.Handle("POST "+prefix+"/monthly-passes", wrap(h.create))
	mux.Handle("POST "+prefix+"/monthly-passes/{id}/toggle-auto-renew", wrap(h.toggleAutoRenew))
}
