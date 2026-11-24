package httpapi

import (
	"net/http"
)

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/status", a.status)
	mux.HandleFunc("POST /api/v1/orders", a.createOrder)
	mux.HandleFunc("GET /api/v1/orders/{id}", a.getOrder)
	// Fallback с тем же путём, но без метода, сохраняет единый JSON-контракт.
	mux.HandleFunc("/api/v1/status", methodNotAllowed("GET, HEAD"))
	mux.HandleFunc("/api/v1/orders", methodNotAllowed("POST"))
	mux.HandleFunc("/api/v1/orders/{id}", methodNotAllowed("GET, HEAD"))
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) { writeError(w, 404, "route not found", "ROUTE_NOT_FOUND") })
	// TODO: после генерации docs подключить Swagger UI здесь.
	return mux
}
func methodNotAllowed(allow string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", allow)
		writeError(w, 405, "method not allowed", "METHOD_NOT_ALLOWED")
	}
}
