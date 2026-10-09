# tg-coffee-shop-bot

Telegram bot for ordering from a coffee shop. Customers browse the menu, build a cart, pay through Telegram Payments and track their order. Store staff process orders through an HTTP API.

## Features

**Customer bot commands**

| Command       | Description                                   |
|---------------|-----------------------------------------------|
| `/start`      | Hint to use `/order`                          |
| `/order`      | Browse categories and products, fill the cart, confirm and pay |
| `/order_info` | Show the active order (status, products, queue position) |
| `/queue`      | Show the current order queue size             |
| `/history`    | Paginated history of past orders              |

After a successful payment, the customer gets a daily queue position, a verification code and a QR code for picking up the order. Customers can cancel an order while it's still allowed. Orders are accepted only while the store is open, according to its schedule.

**Manager HTTP API** (port `8080`, OpenAPI docs at `/docs`)

- `GET /order/next/`: take the next pending order to process
- `PATCH /order/{id}/`: update an order's status; the customer is notified in Telegram

Order lifecycle: `waiting_payment` → `payment_in_progress` → `confirmed` → `in_progress` → `ready` → `completed` (or `canceled` / `rejected`).

## Components

| Binary / service       | Description |
|------------------------|-------------|
| `cmd/bot`              | Telegram bot for customers |
| `cmd/manager`          | HTTP API for store staff |
| `cmd/outboxpoller`     | Polls the outbox tables and sends messages, invoices and payment answers to Telegram |
| `cmd/msgconsumer`      | Alternative to the poller: consumes outbox changes from Kafka (Debezium CDC) |
| `cmd/adminpanel`       | Django admin panel for managing data in the database |

User notifications use the transactional outbox pattern. The bot and the manager write outgoing messages to outbox tables in the same database transaction as the business change. A separate dispatcher (`outboxpoller` or `msgconsumer`) then delivers them to Telegram.

Stack: Go, PostgreSQL, Redis (cart, button callback data, daily order positions), OpenTelemetry + Jaeger, Kafka + Debezium (optional).

## Getting started

### Configuration

Each service reads a YAML config passed with the `-config` flag. Copy the examples and fill in your Telegram bot token and payment provider token:

```bash
cp config/config-bot-example.yaml            config/config-bot.yaml
cp config/config-http-manager-example.yaml   config/config-http-manager.yaml
cp config/config-outbox-poller-example.yaml  config/config-outbox-poller.yaml
cp config/config-msgconsumer-example.yaml    config/config-msgconsumer.yaml
cp config/dbconfig-example.yml               config/dbconfig.yml
```

### Run with Docker Compose

```bash
make compose-up          # postgres, redis, jaeger, migrations, bot, manager, outbox poller
make load-test-data      # optional: load sample store, schedule, currency, categories, products and orders
make compose-down
```

Debezium/Kafka variant (uses `msgconsumer` instead of the poller):

```bash
make debezium-compose-up
make debezium-compose-down
```

Jaeger UI: http://localhost:16686

### Admin panel

```bash
make run-django-admin    # creates a Python venv in ./bin and starts the Django dev server
```

### Kubernetes (minikube)

```bash
make minikube-load-images
make minikube-apply
make minikube-delete
```

## Development

```bash
make build      # build binaries into ./bin
make test       # run tests (repository tests start Postgres in Docker via dockertest)
make lint       # golangci-lint
make fmt        # format code
make generate   # regenerate mocks
make vendor     # tidy and vendor dependencies
```

Database migrations are in `script/db/migrations`.
