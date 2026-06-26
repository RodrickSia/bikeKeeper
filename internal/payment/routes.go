package payment

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, prefix string, mw ...func(http.Handler) http.Handler) {
	wrap := func(hf http.HandlerFunc) http.Handler {
		var handler http.Handler = hf
		for i := len(mw) - 1; i >= 0; i-- {
			handler = mw[i](handler)
		}
		return handler
	}

	mux.Handle("POST "+prefix+"/cards/{cardUID}/deposit", wrap(h.deposit))
	mux.Handle("POST "+prefix+"/cards/{cardUID}/withdraw", wrap(h.withdraw))
	mux.Handle("GET "+prefix+"/cards/{cardUID}/transactions", wrap(h.listByCard))
	mux.Handle("GET "+prefix+"/cards/{cardUID}/balance", wrap(h.getBalance))
}
