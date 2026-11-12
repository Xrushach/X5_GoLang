package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"lesson8/internal/httpapi"
	"lesson8/internal/model"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: httpapi.NewRouter(model.NewStore()),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
