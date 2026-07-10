package main

import (
	"log"
	"net/http"

	"github.com/tobiashenke/go_k8s/internal/items"
)

const redisAddr = "localhost:6379"

func main() {
	repo := items.NewInMemoryItemRepository()
	c := items.NewItemCache(redisAddr)
	service := items.NewItemService(repo, c)
	itemHandler := items.NewItemHandler(service)

	// Setup the server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", itemHandler.HandleGetAll)
	mux.HandleFunc("POST /items", itemHandler.HandleCreate)
	mux.HandleFunc("GET /items/{id}", itemHandler.HandleGetByID)

	// Start the server
	err := http.ListenAndServe(":8087", items.LoggingMiddleware((mux)))
	if err != nil {
		log.Fatal(err)
	}
}
