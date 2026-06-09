package device

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/devices", wrap(h.list, authenticated, staffOnly))
	mux.Handle("GET "+prefix+"/devices/{id}", wrap(h.getByID, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/devices", wrap(h.create, authenticated, staffOnly))
	mux.Handle("PUT "+prefix+"/devices/{id}/status", wrap(h.updateStatus, authenticated, staffOnly))
	mux.Handle("PUT "+prefix+"/devices/{id}/notes", wrap(h.updateNotes, authenticated, staffOnly))
	mux.Handle("GET "+prefix+"/devices/alerts", wrap(h.getAlerts, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/devices/{id}/fault", wrap(h.reportFault, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/devices/alerts/{alertId}/resolve", wrap(h.resolveAlert, authenticated, staffOnly))
}
