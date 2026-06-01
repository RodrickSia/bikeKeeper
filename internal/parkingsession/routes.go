package parkingsession

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/sessions/checkin", wrap(h.checkIn))
	mux.Handle("POST "+prefix+"/sessions/{id}/checkout", wrap(h.checkOut))
	mux.Handle("GET "+prefix+"/sessions/{id}", wrap(h.getByID))
	mux.Handle("GET "+prefix+"/sessions/card/{cardUID}", wrap(h.listByCard))
	mux.Handle("DELETE "+prefix+"/sessions/{id}", wrap(h.delete))
}
