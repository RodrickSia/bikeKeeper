package parkinglot

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, facultyOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("GET "+prefix+"/parking-lots", wrap(h.list, authenticated))
	mux.Handle("GET "+prefix+"/parking-lots/{id}", wrap(h.getByID, authenticated))
	mux.Handle("POST "+prefix+"/parking-lots", wrap(h.create, authenticated, facultyOnly))
	mux.Handle("PUT "+prefix+"/parking-lots/{id}", wrap(h.update, authenticated, facultyOnly))
	mux.Handle("DELETE "+prefix+"/parking-lots/{id}", wrap(h.delete, authenticated, facultyOnly))
}
