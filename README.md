# Orders Worker

Internal service that consumes order commands from Kafka, stores them in
Postgres, and serves order reads over gRPC. It has **no public HTTP API** — the
only entry points are the Kafka consumer and the gRPC server. The public
`test-api-service` talks to it; nothing reaches it from the internet.

```
test-api-service ──POST/PATCH/DELETE──▶ Kafka ──▶ worker ──▶ Postgres
test-api-service ──────────GET (gRPC)──────────▶ worker ──▶ Postgres
```

Kafka is started by `test-api-service`; this service starts Postgres. Each
service is independent with its own `compose.dev.yml` and is brought up
separately (`make run` in each). Both use host networking (`network_mode: host`),
so they share the host's network namespace and reach each other over `localhost`
— no shared network to create.

## Layers

The orders domain sits in a funnel `services/worker/orders` (mirroring the API
service's `services/api/mobile/orders`):

- **handler** — `src/services/worker/orders/consumer.go` (Kafka) and `handler.go` (gRPC)
- **service** — `src/services/worker/orders/services.go` (validation, business rules)
- **repository** — `src/services/worker/orders/repositories.go` (Postgres via GORM)

## Prerequisites

- Docker, Go 1.26.
- gRPC codegen runs through the Go toolchain (buf). Atlas runs natively — install
  it once from https://atlasgo.io (e.g. `scoop install atlas`). Running Atlas in
  Docker is avoided on purpose: bind-mounting a non-ASCII project path (Cyrillic /
  spaces) breaks under GnuWin32 make.

## Run the full stack

Start each service separately (run `make run` in `test-api-service` for
Kafka + the API, and here for the worker + Postgres). Both use host networking,
so they end up reachable over `localhost`.

```powershell
make run        # builds the binary, starts worker app + Postgres
make migrate-up # applies migrations to the database
make down       # stops the stack
```

Migrations need a running Postgres; the consumer retries until Kafka and its
topic are available, so start order is not strict.

## Local development

```powershell
make deps       # starts only Postgres
make migrate-up
make local-run  # go run -race ./src
```

## Migrations (Atlas + Gorm loader)

The schema source of truth is the Gorm models in `src/core/models`. The Atlas
Gorm loader (`src/core/db/atlas_loader`) turns them into SQL; Atlas diffs that
against `migrations/` to produce versioned migrations. Atlas runs in Docker.

```powershell
make makemigrations name=add_orders   # generate a new migration from the models
make migrate-up                        # apply pending migrations
make migrate-status                    # show applied / pending migrations
make migrate-down                      # roll back the last migration
```

To change the schema: edit a model under `src/core/models`, register it in
`src/core/models/registry.go`, then run `make makemigrations name=...`.

## gRPC contract

The contract lives in `proto/orders/orders.proto` (no API versioning) and is
shared with `test-api-service`. Regenerate the Go code after editing the proto:

```powershell
make proto-tools   # installs buf + protoc-gen-go(-grpc) via the Go toolchain (once)
make proto-gen     # regenerates src/proto/**
```

## Configuration

```text
KAFKA_BROKER_ADDRESSES=localhost:9092       # where Kafka is reachable
KAFKA_ORDERS_TOPIC=orders                   # topic to consume
KAFKA_ORDERS_CONSUMER_GROUP=orders-worker   # consumer group id
DATABASE_URL=postgres://orders:orders@localhost:5432/orders?sslmode=disable
GRPC_LISTEN_ADDRESS=:50051                  # gRPC server bind address
```

## Checks

```powershell
make check   # go mod tidy / fmt / vet / test
```
