# 3-Day Learning Plan: Go + Distributed Systems + K8s Operator

## Overview

Each day builds on the last. The operator only appears at the very end of day 3.

**Your mental model bridge (ops → dev):**

| You know | Equivalent here |
|---|---|
| Terraform resource | Go struct with typed fields |
| Terraform plan/apply | Operator reconcile loop |
| Helm values.yaml | Custom Resource spec |
| CRD schema | Go types → `make generate` |
| k8s deployment manifest | What your operator creates programmatically |

---

## Theory Reference

Read each topic before the day it's marked as required. Each entry is intentionally short — the goal is the mental model, not exhaustive detail.

| # | Topic | One-line explanation | Required before | Done |
|---|---|---|---|---|
| T1 | HTTP fundamentals | Methods, status codes, headers, idempotency (GET/PUT safe to retry, POST is not) | Day 1 | [x] |
| T2 | REST design principles | Resources are nouns, verbs are HTTP methods, state lives on the server not the client | Day 1 | [x] |
| T3 | API error design | Errors need a consistent shape — RFC 7807 "Problem Details" is the standard (`type`, `title`, `status`, `detail`) | Day 1 | [x] |
| T4 | Dependency Injection | Pass dependencies in from outside rather than constructing them inside — enables swapping implementations without changing callers | Day 1 | [x] |
| T5 | Repository pattern | Abstract all data access behind an interface (`type ItemRepository interface { Get, Save, Delete }`) — business logic never touches storage directly | Day 1 | [x] |
| T6 | Service layer pattern | HTTP handler → Service (business logic) → Repository (storage) — each layer knows nothing about the layers below it | Day 1 | [x] |
| T7 | Middleware pattern | A function that wraps a handler to add cross-cutting behaviour (logging, auth, rate limiting) without touching the handler itself | Day 1 optional | [x] |
| T8 | Caching strategies | Cache-aside (read): check cache first, on miss load from DB and write to cache. Write-through: write to DB and cache together. Write-behind: write to cache, flush to DB async | Day 2 | [x] |
| T9 | TTL, LRU eviction, cache stampede | TTL: entry expires after fixed time. LRU: evict least-recently-used when full. Stampede: all requests miss simultaneously → all hit DB → overload | Day 2 | [x] |
| T10 | Messaging delivery guarantees | At-most-once: fire and forget, may lose. At-least-once: retried until acked, may duplicate. Exactly-once: guaranteed once, expensive and often a lie in practice | Day 2 | [x] |
| T11 | Idempotency | Same operation called N times produces the same result as calling it once — required for at-least-once consumers and safe retries | Day 2 | [x] |
| T12 | Dead letter queue (DLQ) | Messages that fail processing repeatedly are moved to a DLQ instead of blocking the main queue — lets you inspect and replay them separately | Day 2 optional | [x] |
| T13 | Database fundamentals | Index = pre-sorted lookup structure (fast reads, slower writes). N+1 problem = one query to list N records + N queries for each record's relations. Connection pool = reuse connections instead of opening one per request | Day 2 | [x] |
| T14 | ORM vs raw SQL | ORM (e.g. GORM): maps structs to tables, handles relations, migration — convenient but hides the SQL and can produce bad queries. Raw SQL (`database/sql`, `sqlx`): explicit, fast, verbose. Go culture leans toward `sqlx` or query builders over full ORMs | Day 2 | [x] |
| T15 | Backpressure | When a consumer is slower than a producer, the queue grows unboundedly — systems need to either slow the producer, drop messages, or scale the consumer | Day 3 | [x] |
| T16 | CAP theorem | A distributed system can guarantee at most two of: Consistency, Availability, Partition tolerance — partitions always happen, so every system trades off C vs A | Day 3 | [x] |
| T17 | ACID vs BASE | ACID (relational DBs): Atomic, Consistent, Isolated, Durable. BASE (distributed systems): Basically Available, Soft state, Eventually consistent — explains why caches and queues behave differently than a DB | Day 3 | [x] |
| T18 | Observability pillars | Logs: what happened (structured, queryable). Metrics: how much / how fast (counters, gauges, histograms). Traces: where time was spent across service boundaries. Different tools for different questions | Day 4 | [x] |
| T19 | Circuit breaker | When a dependency fails repeatedly, stop calling it for a cooldown period (open state) — fail fast instead of piling up slow failures that exhaust threads/goroutines | Day 5 | [ ] |
| T20 | Retry with exponential backoff + jitter | On failure, wait 2^n seconds before retrying — jitter (random offset) prevents all retriers hitting the service simultaneously (thundering herd) | Day 5 | [x] |
| T21 | Fail open vs fail closed | Fail open: allow requests through when a dependency is unavailable (availability over safety). Fail closed: reject requests when a dependency is down (safety over availability). Rate limiting and auth caching fail open; payments and auth validation fail closed | Day 6 | [ ] |
| T22 | Reverse proxy & Ingress | A reverse proxy sits in front of services, routing external traffic by hostname/path, handling TLS termination, load balancing. In k8s, an Ingress resource configures this — backed by nginx, Traefik, or cloud LBs | Day 7 | [ ] |

---

## Day 1 — Go Fundamentals

**Goal:** You can write and run a working HTTP server in Go from scratch.

### ~~Morning (2-3h): Language basics~~ ✅ Done
- ~~[A Tour of Go](https://go.dev/tour)~~ — completed via [w3schools.com/go](https://www.w3schools.com/go/index.php)
- Key insight: Go structs + interfaces are like Terraform resource schemas — typed, explicit, composed

### ~~Afternoon (2-3h): Build a small HTTP API~~ ✅ Done
- Write a REST API with the standard library (`net/http`) — no frameworks
- Endpoints: `GET /items`, `POST /items`, `GET /items/{id}`
- Use `encoding/json` for marshaling
- This solidifies: structs, error handling, handlers, JSON

### ~~Evening (1h): Go modules & tooling~~ ✅ Done
- `go mod init`, `go get`, `go run`, `go build`, `go test`
- Think of `go.mod` as your `versions.tf`

### ~~Optional (if ahead of schedule)~~ ✅ Done
- ~~Add middleware: a simple request logger that wraps your handlers~~ ✅
- ~~Try goroutines: process items concurrently using `sync.WaitGroup`~~ ✅
- ~~Read [Effective Go](https://go.dev/doc/effective_go) sections on concurrency — channels, goroutines, CSP model~~ ✅

---

## Day 2 — Messaging + Caching + Distributed Patterns

**Goal:** Running services that communicate via events and serve cached reads.

Key mental model: **messaging decouples producers from consumers; caching decouples hot reads from slow storage**.

### ~~Morning (2h): Redis / Caching~~ ✅ Done
### ~~Midday (2-3h): Messaging with NATS~~ ✅ Done
### ~~Afternoon (2h): Wire it together~~ ✅ Done

### ~~Optional (if ahead of schedule)~~ ✅ Done
- ~~Implement cache invalidation: delete the cache key on `DELETE /items/{id}`~~ ✅
- ~~Add structured logging with `log/slog` (stdlib, Go 1.21+) and include a trace/correlation ID in every log line~~ ✅
- ~~Replace the in-memory map with a real SQLite store using raw SQL (`database/sql`)~~ ✅
- ~~Migrated SQLite repository from raw SQL to GORM ORM~~ ✅

---

## Day 3 — Kubernetes Operator

**Goal:** A running operator that responds to custom resources.

An operator is just: watch for custom resources → reconcile state.

### ~~Morning (2h): Controller concepts~~ ✅ Done
- ~~Read the [operator pattern docs](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)~~ — covered in theory
- Key mental model: a controller is a `for` loop that watches k8s API events and drives actual state toward desired state — the same reconcile loop Terraform does, but continuously running inside the cluster instead of on your laptop
- ~~Install tools:~~
  ```bash
  brew install kubebuilder
  brew install kind
  kind create cluster
  ```
- Learned concepts: CRD, operator, controller, reconciler, manager, scheme, watch stream, work queue, RBAC, controller-runtime, Kustomize, leader election, control plane

### ~~Midday (3h): Scaffold and implement a basic operator~~ ✅ Done

```bash
kubebuilder init --domain henkebyte.dev --repo github.com/tobiashenke/go_k8s/operator
kubebuilder create api --group apps --version v1alpha1 --kind Widget
```

**Operator goal:** A `Widget` custom resource that creates a `ConfigMap` with the widget's spec fields as data. That's it — no complex business logic. The point is to understand the reconcile loop.

1. ~~Edit the generated `spec` struct in `api/v1alpha1/widget_types.go`~~ ✅
2. ~~Implement `Reconcile()` in `internal/controller/widget_controller.go`: get the Widget, create/update a ConfigMap~~ ✅
3. ~~`make install` to install the CRD, `make run` to run the controller locally against your cluster~~ ✅

### Afternoon (1-2h): Connect your Day 2 work
- In the reconciler, also publish a NATS event when a Widget is created/updated
- This ties everything together: Go + messaging + k8s operator

### ~~Optional (if ahead of schedule)~~ ✅ Done
- ~~Add a `status` subresource to your Widget CRD and update it in the reconciler (`updateStatus`)~~ ✅
- Handle deletion with a finalizer — prevents the resource being deleted before cleanup runs
- ~~Deploy the operator into your kind cluster as a real Pod: `make docker-build`, `make deploy`~~ ✅

---

## Day 4 (Optional) — Production-Grade API

**Goal:** Turn your Day 1-2 API into something that looks like a real service.

- ~~Swap `net/http` for [chi](https://github.com/go-chi/chi) — learn routing, middleware, request binding~~ ✅
- ~~Add proper validation with `github.com/go-playground/validator/v10`~~ ✅
- ~~Add structured logging (`log/slog`) — completed in Day 2~~ ✅
- ~~Add Prometheus metrics: request count, latency histograms — use `github.com/prometheus/client_golang`~~ ✅
- ~~Write a `docker-compose.yml` that wires up your API, Redis and NATS together~~ ✅
- ~~Add a health check endpoint (`GET /healthz`) and a readiness endpoint (`GET /readyz`) — standard for k8s deployments~~ ✅
- ~~Deploy the full stack to your kind cluster using a Helm chart you write yourself~~ ✅

---

## Day 5 (Optional) — Resilience Patterns

**Goal:** Idempotency, rate limiting, distributed tracing, Lua fix. (~5-6h)

- ~~**Idempotency:** Make your `POST /items` idempotent using a client-supplied `Idempotency-Key` header cached in Redis~~ ✅ (~2h)
- ~~**Rate limiting:** Add a Redis-backed sliding-window rate limiter to your API~~ ✅ (~2-3h)
- ~~**Distributed tracing:** Instrument your operator and API with OpenTelemetry traces — see a request flow across services in Jaeger~~ ✅ (~1.5h)
- **Rate limiter race condition fix:** Implement atomic sliding window rate limiting using a Redis Lua script (~45min)

---

## Day 6 (Optional) — Security & Reliability

**Goal:** Authentication, circuit breaker, Postgres. (~4-5h)

- **Authentication:** Add JWT-based authentication middleware to your API — validate tokens on protected endpoints (~1.5h)
- **Circuit breaker:** Implement a basic circuit breaker around your cache calls using `github.com/sony/gobreaker` (~45min)
- **Postgres:** Add a Postgres service to docker-compose, swap SQLite repository for a `PostgresItemRepository` using `database/sql` — same interface, different driver (`github.com/lib/pq`), learn about connection strings and connection pooling (~1h)

---

## Day 7 (Optional) — Messaging, Data & Operators

**Goal:** Kafka, operator conditions, proxy, reading. (~5-6h)

- **Kafka deep dive:** Replace NATS with Kafka, implement consumer groups, understand partition assignment and offset management (~1.5h)
- **Operator status conditions:** Implement proper `metav1.Condition` status conditions on your Widget — this is the standard pattern used by all production operators (cert-manager, crossplane, etc.) (~45min)
- **Reverse proxy:** Deploy an nginx or Traefik reverse proxy in front of your API in the kind cluster — learn about Ingress resources, path routing, TLS termination (~2h)
- **Read:** [Designing Distributed Systems](https://www.oreilly.com/library/view/designing-distributed-systems/9781491983638/) ch. 1-4 (free online) — patterns like sidecar, ambassador, adapter map directly to k8s (~1h)

---

## Day 8 (Optional) — LLM Agent in Kubernetes

**Goal:** Deploy an AI agent to the cluster that can operate on your data. (~5-7h)

- **LLM-powered API agent:** Build a Go service using the Claude API that accepts natural language tasks (e.g. "show me all items", "delete item 1", "create an item called widget") and translates them to calls against your existing HTTP API — deploy it to kind alongside your API
- **Operator-driven agent:** Alternatively, define an `AgentTask` CRD — users apply a YAML describing a task, your operator reconciles it by calling the Claude API and executing the result against the database
- **Tool use pattern:** Implement Claude tool use — give the model tools like `get_items`, `create_item`, `delete_item` backed by your existing service layer
- Ideas for useful agents in your setup: natural language query interface, automated data cleanup, anomaly detection on items, NATS event summariser

---

## Key Resources

| Topic | Resource |
|---|---|
| Go language | [go.dev/tour](https://go.dev/tour) |
| Effective Go | [go.dev/doc/effective_go](https://go.dev/doc/effective_go) |
| Go HTTP | [pkg.go.dev/net/http](https://pkg.go.dev/net/http) |
| Redis (go-redis) | [redis.uptrace.dev](https://redis.uptrace.dev) |
| NATS | [docs.nats.io](https://docs.nats.io) |
| kubebuilder | [book.kubebuilder.io](https://book.kubebuilder.io) — ch. 1-3 only |
| Operator pattern | [kubernetes.io/docs/concepts/extend-kubernetes/operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) |

---

## Where to start right now

[go.dev/tour](https://go.dev/tour) — get through it today. Everything else unlocks from there.


## Guides and Roadmaps

Golang: https://roadmap.sh/golang
Backend: https://roadmap.sh/backend
