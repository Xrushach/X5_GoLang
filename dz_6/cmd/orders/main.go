package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"example.com/lesson8/internal/httpapi"
	"example.com/lesson8/internal/repository"
	"example.com/lesson8/internal/service"
)

func main() {
	router := flag.String("router", "mux", "mux или chi")
	demo := flag.Bool("demo", false, "включить GET /demo/panic")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	api := httpapi.New(service.New(repository.New()), logger)
	var handler http.Handler
	switch *router {
	case "mux":
		handler = api.Mux(*demo)
	case "chi":
		handler = api.Chi(*demo)
	default:
		logger.Error("unknown router", "router", *router)
		os.Exit(1)
	}
	server := &http.Server{
		Addr: "127.0.0.1:8080", Handler: handler,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	logger.Info("server starting", "addr", server.Addr, "router", *router)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
