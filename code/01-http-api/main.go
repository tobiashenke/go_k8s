package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/tobiashenke/go_k8s/internal/items"
)

const redisAddr = "localhost:6379"
const natsUrl = "nats://localhost:4222"

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Business logic
	repo := items.NewInMemoryItemRepository()
	c := items.NewItemCache(redisAddr)
	p, err := items.NewItemPublisher(natsUrl)
	if err != nil {
		slog.Error("failed to start the NATS publisher", "error", err)
		os.Exit(1)
	}
	service := items.NewItemService(repo, c, p)
	itemHandler := items.NewItemHandler(service)

	// Setup the server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", itemHandler.HandleGetAll)
	mux.HandleFunc("POST /items", itemHandler.HandleCreate)
	mux.HandleFunc("GET /items/{id}", itemHandler.HandleGetByID)
	mux.HandleFunc("DELETE /items/{id}", itemHandler.HandleDelete)

	// Start the server
	err = http.ListenAndServe(":8087", items.LoggingMiddleware((mux)))
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
