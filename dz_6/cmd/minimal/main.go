// Первый показ: Server, ServeMux, Handler и HandlerFunc.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type orderPreview struct{}

func (orderPreview) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Откуда берётся id? Существует ли уже такой заказ в хранилище?
	fmt.Fprintf(w, "Маршрут заказа, id=%s\n", r.PathValue("id"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Какой статус получится без WriteHeader?
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.Handle("GET /orders/{id}", orderPreview{})

	// Что изменится, если передать nil вместо mux?
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	log.Printf("Сервер слушает http://%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
