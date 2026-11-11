package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *API) Mux(demo bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("POST /orders", a.createOrder)
	mux.HandleFunc("GET /orders/{id}", a.getOrder)
	if demo {
		mux.HandleFunc("GET /demo/panic", demoPanic)
	}
	// В каком порядке функции начнут и закончат работу?
	return a.requestID(a.logging(a.recovery(mux)))
}

func (a *API) Chi(demo bool) http.Handler {
	router := chi.NewRouter()
	router.Use(a.requestID, a.logging, a.recovery)
	router.Get("/health", a.health)
	router.Head("/health", a.health)
	router.Route("/orders", func(r chi.Router) {
		r.Post("/", a.createOrder)
		r.Get("/{id}", a.chiGetOrder)
		r.Head("/{id}", a.chiGetOrder)
	})
	if demo {
		router.Get("/demo/panic", demoPanic)
	}
	return router
}

func (a *API) chiGetOrder(w http.ResponseWriter, r *http.Request) {
	// Адаптер параметра позволяет использовать тот же обработчик для обоих роутеров.
	r.SetPathValue("id", chi.URLParam(r, "id"))
	a.getOrder(w, r)
}

func demoPanic(w http.ResponseWriter, r *http.Request) {
	panic("демонстрационная паника")
}
