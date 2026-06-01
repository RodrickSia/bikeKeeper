package member

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/members", wrap(h.create))
	mux.Handle("GET "+prefix+"/members", wrap(h.list))
	mux.Handle("GET "+prefix+"/members/{id}", wrap(h.getByID))
	mux.Handle("GET "+prefix+"/members/student/{studentID}", wrap(h.getByStudentID))
	mux.Handle("PUT "+prefix+"/members/{id}", wrap(h.update))
	mux.Handle("DELETE "+prefix+"/members/{id}", wrap(h.delete))
}
