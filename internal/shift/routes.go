package shift

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/shifts", wrap(h.list, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/shifts", wrap(h.create, authenticated, staffOnly))
	mux.Handle("GET "+prefix+"/shifts/{id}", wrap(h.getByID, authenticated, staffOnly))
	mux.Handle("PUT "+prefix+"/shifts/{id}/status", wrap(h.updateStatus, authenticated, staffOnly))
	mux.Handle("PUT "+prefix+"/shifts/{id}/notes", wrap(h.updateNotes, authenticated, staffOnly))
	mux.Handle("DELETE "+prefix+"/shifts/{id}", wrap(h.delete, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/shifts/{id}/assign", wrap(h.assignStaff, authenticated, staffOnly))
	mux.Handle("POST "+prefix+"/shifts/{id}/unassign", wrap(h.removeStaff, authenticated, staffOnly))
}
