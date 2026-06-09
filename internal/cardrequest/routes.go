package cardrequest

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, facultyOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	// Student routes
	mux.Handle("POST "+prefix+"/card-requests", wrap(h.submit, authenticated))
	mux.Handle("GET "+prefix+"/card-requests/me", wrap(h.listMine, authenticated))

	// Admin/Faculty routes
	mux.Handle("GET "+prefix+"/card-requests", wrap(h.list, authenticated, facultyOnly))
	mux.Handle("POST "+prefix+"/card-requests/{id}/approve", wrap(h.approve, authenticated, facultyOnly))
	mux.Handle("POST "+prefix+"/card-requests/{id}/reject", wrap(h.reject, authenticated, facultyOnly))
	mux.Handle("POST "+prefix+"/card-requests/{id}/block", wrap(h.block, authenticated, facultyOnly))
}
