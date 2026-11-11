package main

import (
	"log"
	"net/http"
	"time"
)

func health(w http.ResponseWriter, r *http.Request) {
	// Как задать формат и статус ответа?
	http.Error(w, "TODO: health", http.StatusNotImplemented)
}

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: newRouter(NewStore()),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
