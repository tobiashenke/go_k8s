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
var postgresDSN = os.Getenv("POSTGRES_DSN")

func main() {
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	if natsUrl == "" {
		natsUrl = "nats://localhost:4222"
	}
	if postgresDSN == "" {
		postgresDSN = "host=localhost user=admin password=secret dbname=items port=5432 sslmode=disable"
	}
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Auth Handler

	// OTEL tracer
	shutdownTracer, err := items.InitTracer(context.Background(), "http-api")
	if err != nil {
		slog.Error("failed to initialize the tracer", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer()

	// Business logic
	repo, err := items.NewPostgresItemRepository(postgresDSN)
	if err != nil {
		slog.Error("failed to connect to Postgres database", "error", err)
		os.Exit(1)
	}
	c, err := items.NewItemCache(redisAddr)
	if err != nil {
		slog.Error("failed to start the redis cache", "error", err)
		os.Exit(1)
	}
	p, err := items.NewItemPublisher(natsUrl)
	if err != nil {
		slog.Error("failed to start the NATS publisher", "error", err)
		os.Exit(1)
	}
	service := items.NewItemService(repo, c, p)
	itemHandler := items.NewItemHandler(service)
	a := items.AuthHandler{}

	// Setup the server
	r := chi.NewRouter()
	// Register middleware
	r.Use(items.LoggingMiddleware)
	r.Use(items.MetricsMiddleware)
	r.Use(items.RateLimitMiddleware(c, 100, 60*time.Second))
	// Register route with middleware
	r.Route("/items", func(r chi.Router) {
		r.Use(items.AuthMiddleWare)
		r.Use(items.IdempotencyMiddleware(c))
		r.Post("/", itemHandler.HandleCreate)
	})
	// Register routes
	r.Post("/login", a.HandleLogin)
	r.Group(func(r chi.Router) {
		r.Use(items.AuthMiddleWare)
		r.Get("/items", itemHandler.HandleGetAll)
		r.Get("/items/{id}", itemHandler.HandleGetByID)
		r.Delete("/items/{id}", itemHandler.HandleDelete)
	})
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
