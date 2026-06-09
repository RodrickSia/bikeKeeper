package vehicle

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, authenticated func(http.Handler) http.Handler, staffOnly func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
		var handler http.Handler = hf
		for i := len(mws) - 1; i >= 0; i-- {
			handler = mws[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/vehicles", wrap(h.add, authenticated))
	mux.Handle("GET "+prefix+"/vehicles/me", wrap(h.myVehicles, authenticated))
	mux.Handle("GET "+prefix+"/vehicles", wrap(h.list, authenticated, staffOnly))
	mux.Handle("GET "+prefix+"/vehicles/plate/{plate}", wrap(h.findByPlate, authenticated, staffOnly))
	mux.Handle("DELETE "+prefix+"/vehicles/{id}", wrap(h.remove, authenticated))
}
