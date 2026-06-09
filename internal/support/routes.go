package support

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	// Student routes
	mux.Handle("POST "+prefix+"/support/tickets", wrap(h.create, authenticated))
	mux.Handle("GET "+prefix+"/support/tickets/me", wrap(h.listMine, authenticated))

	// Authenticated routes (either user or staff/admin can view/respond to a specific ticket)
	mux.Handle("GET "+prefix+"/support/tickets/{id}", wrap(h.getByID, authenticated))
	mux.Handle("POST "+prefix+"/support/tickets/{id}/responses", wrap(h.addResponse, authenticated))

	// Staff/Admin/Faculty routes
	mux.Handle("GET "+prefix+"/support/tickets", wrap(h.listAll, authenticated, staffOnly))
	mux.Handle("PUT "+prefix+"/support/tickets/{id}/status", wrap(h.updateStatus, authenticated, staffOnly))
}
