package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/tobiashenke/go_k8s/internal/items"
)

const redisAddr = "localhost:6379"
const natsUrl = "nats://localhost:4222"

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Business logic
	repo, err := items.NewSQLiteItemRepository()
	if err != nil {
		slog.Error("failed to open SQLite database", "error", err)
		os.Exit(1)
	}
	c := items.NewItemCache(redisAddr)
	p, err := items.NewItemPublisher(natsUrl)
	if err != nil {
		slog.Error("failed to start the NATS publisher", "error", err)
		os.Exit(1)
	}
	service := items.NewItemService(repo, c, p)
	itemHandler := items.NewItemHandler(service)

	// Setup the server
	r := chi.NewRouter()
	// Register middleware
	r.Use(items.LoggingMiddleware)
	// Register routes
	r.Get("/items", itemHandler.HandleGetAll)
	r.Post("/items", itemHandler.HandleCreate)
	r.Get("/items/{id}", itemHandler.HandleGetByID)
	r.Delete("/items/{id}", itemHandler.HandleDelete)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Start the server
	err = http.ListenAndServe(":8087", r)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
