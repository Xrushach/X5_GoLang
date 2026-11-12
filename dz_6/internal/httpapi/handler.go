package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"lesson8/internal/model"

	"github.com/go-chi/chi/v5"
)

/*
func newRouter(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /orders", createOrder(store))
	return mux
}
*/

func NewRouter(store *model.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(RequestID)
	r.Use(Logger)
	r.Use(Recovery)
	r.Get("/health", health)
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", createOrder(store))
		r.Get("/{id}", getOrder(store))
	})
	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func createOrder(store *model.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Прочитайте customer из JSON и проверьте имя до вызова store.Create.
		// Как завершить обработку сразу после ошибки?
		var data struct {
			Customer string `json:"customer"`
		}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&data); err != nil {
			writeJSON(w, 400, map[string]string{"error": "Invalid JSON"})
			return
		}
		customer := strings.TrimSpace(data.Customer)
		if customer == "" {
			writeJSON(w, 400, map[string]string{"error": "Customer is empty"})
			return
		}
		order := store.Create(customer)
		w.Header().Set("Location", fmt.Sprintf("/orders/%d", order.ID))
		writeJSON(w, 201, order)
	}
}

func getOrder(store *model.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil || id <= 0 {
			writeJSON(w, 400, map[string]string{"error": "Invalid ID"})
			return
		}
		order, ok := store.ByID(id)
		if !ok {
			writeJSON(w, 404, map[string]string{"error": "Order not found"})
			return
		}
		writeJSON(w, 200, order)
	}
}
