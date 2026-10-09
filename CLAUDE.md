# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Telegram bot for a coffee shop: customers browse products, build a cart, pay via Telegram invoices, and track orders; a manager HTTP API moves orders through statuses.

## Commands

```bash
make build        # builds all Go binaries into ./bin (uses -mod=vendor)
make test         # go test ./...
make lint         # golangci-lint v2 (auto-installed into tools/bin) with .golangci.yml
make fmt          # gofmt + goimports via golangci-lint
make generate     # go generate ./... (mockgen is a go.mod `tool` dependency)
make vendor       # go mod tidy && go mod vendor — run after changing dependencies
make compose-up   # full stack: postgres, redis, jaeger, sql-migrate, bot, httpmanager, outboxpoller
make debezium-compose-up  # alternative stack using Debezium CDC -> Kafka -> msgconsumer
make load-test-data       # load script/db/dataset/test_data.sql into local postgres
```

Run a single test: `go test ./internal/adapter/repository/postgres/pgorder/ -run TestPgOrderSuit/TestName`

- Dependencies are vendored (`vendor/`); builds use `-mod=vendor`, so `make vendor` is required after editing `go.mod`.
- Postgres repository tests use `ory/dockertest` to spin up a real Postgres container and apply `script/db/migrations` via `sql-migrate` — Docker must be running.
- The linter enables nearly all linters (`default: all`); expect strict rules (funlen, gochecknoglobals, varnamelen, wrapcheck, etc.). Use targeted `//nolint:<linter>` only when justified, as existing code does.
- CI (GitHub Actions) runs `make build`, `make test`, and golangci-lint on Go 1.26.

## Configuration

Each binary loads YAML via `configor` from the `-config` flag (default `config/config.yaml`). Example configs live in `config/*-example.yaml`; real ones (`config/config-bot.yaml`, etc.) are gitignored and mounted by docker-compose. Config structs implement `application.Configer` (log level, service name, tracing endpoint).

## Architecture

Hexagonal / ports-and-adapters layout with several binaries in `cmd/` sharing `internal/`:

- `cmd/bot` — Telegram bot (uses `github.com/Mikhalevich/tgbot` wrapper over `go-telegram/bot`). Routes in `cmd/bot/internal/app/routes.go`, handlers in `tghandler/`.
- `cmd/manager` — HTTP API for store staff (huma v2): fetch next pending order, update order status.
- `cmd/outboxpoller` — polls outbox tables and dispatches messages/invoices/payment answers to Telegram.
- `cmd/msgconsumer` — alternative outbox dispatcher: consumes Debezium CDC events from Kafka.
- `cmd/adminpanel` — Django admin (Python) over the same DB, run with `make run-django-admin`.

Every Go binary follows the same shape: `main.go` calls `application.Run(&cfg, fn)` (config, logrus logger, OTel tracing, signal handling), and `internal/setup/setup.go` does all manual dependency wiring (no DI framework).

`internal/` layers:

- `domain/port/` — domain types and value objects (order, product, store, msginfo, button, etc.) and `perror` (typed domain errors; check with `perror.IsType`).
- `domain/service/` — reusable domain services (`ordersvc`, `cartsvc`, `productsvc`, `notificationsvc`, `outbox/outboxsvc`, ...).
- `domain/usecase/` — application use cases split by actor (`customer/...`, `manager/...`). Each use case package declares the narrow interfaces it consumes (`OrderService`, `NotificationService`, ...) — consumer-side interfaces. Services assert conformance with `var _ usecase.Iface = (*Service)(nil)` blocks.
- `adapter/` — implementations: Postgres repositories (`repository/postgres/pg*`, sqlx + pgx, with private `internal/model` row structs), Redis-backed cart / daily position counter / button storage, Telegram `messagesender`, QR and verification code generators.
- `infra/` — application bootstrap, logger (logrus, carried in context via `logger.FromContext`), tracing (OpenTelemetry/Jaeger).

Key cross-cutting patterns:

- **Transactions via context**: `transaction.Transaction(ctx, fn)` stores the tx in the context; repositories pick it up automatically and nested calls reuse the outer tx.
- **Transactional outbox**: in `bot` and `manager`, `notificationsvc` is constructed with `pgoutbox` as its sender, so user-facing Telegram messages/invoices are written to outbox tables inside the business transaction rather than sent directly. `outboxpoller` (or Debezium → Kafka → `msgconsumer`) later delivers them via `messagesender`.
- **Inline button callbacks**: button payloads are stored in Redis (`buttonrespository`) and referenced by ID in Telegram callback data.

DB schema lives in `script/db/migrations/*.sql` (applied by the `sql-migrate` container in compose and by repository tests). Kubernetes manifests for minikube are in `script/k8s/minikube`.
