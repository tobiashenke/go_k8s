package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tobiashenke/go_k8s/internal/items"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var redisAddr = os.Getenv("REDIS_ADDR")
var natsUrl = os.Getenv("NATS_URL")

func main() {
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	if natsUrl == "" {
		natsUrl = "nats://localhost:4222"
	}
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// OTEL tracer
	shutdownTracer, err := items.InitTracer(context.Background(), "http-api")
	if err != nil {
		slog.Error("failed to initialize the tracer", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer()

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
	r.Use(items.MetricsMiddleware)
	r.Use(items.RateLimitMiddleware(c, 100, 60*time.Second))
	// Register route with middleware
	r.Route("/items", func(r chi.Router) {
		r.Use(items.IdempotencyMiddleware(c))
		r.Post("/", itemHandler.HandleCreate)
	})
	// Register routes
	r.Get("/items", itemHandler.HandleGetAll)
	r.Get("/items/{id}", itemHandler.HandleGetByID)
	r.Delete("/items/{id}", itemHandler.HandleDelete)
	// Not required anymore due to middleware POST route
	// r.Post("/items", itemHandler.HandleCreate)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/metrics", http.HandlerFunc(promhttp.Handler().ServeHTTP))

	// Start the server, wrapped chi router into OTEL trace
	err = http.ListenAndServe(":8087", otelhttp.NewHandler(r, "http-api"))
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
