![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ivmerk/46bf0f7a10ffaf0e26dc0646ad729a46/raw/coverage.json)

# goshop

A small spare-parts ordering platform built as a set of Go microservices. Customers create an order for a list of parts, the system prices it from the parts catalog, and the order is then paid through a payment service.

## Services

| Service       | Transport      | Storage    | Responsibility                                           |
| ------------- | -------------- | ---------- | -------------------------------------------------------- |
| **inventory** | gRPC           | MongoDB    | Parts catalog: get a part, list parts with filters       |
| **order**     | REST (OpenAPI) | PostgreSQL | Create, read and pay orders; calls inventory and payment |
| **payment**   | gRPC           | none       | Processes a payment and returns a transaction ID         |

```
client ──REST──▶ order ──gRPC──▶ inventory ──▶ MongoDB
                  │
                  ├──gRPC──▶ payment
                  └──▶ PostgreSQL
```

Order flow: `POST /api/v1/orders` validates the parts through inventory and stores the order as `PENDING_PAYMENT`. `POST /api/v1/orders/{uuid}/pay` calls payment and moves the order to `PAID`.

## Repository layout

```
.
├── inventory/      # gRPC service, MongoDB
├── order/          # REST service, PostgreSQL, gRPC clients
├── payment/        # gRPC service
├── platform/       # shared library: logger, closer, migrator, gRPC health, testcontainers helpers
├── shared/         # protobuf and OpenAPI definitions, generated code
└── deploy/
    ├── compose/    # Docker Compose per service
    ├── docker/     # Dockerfiles
    └── env/        # env templates and generator
```

Each service is a separate Go module in one workspace (`go.work`). Inside a service the layers are `api → service → repository`, wired together in `internal/app` through a small dependency-injection container.

## Tech stack

- Go 1.27, Go workspaces
- gRPC and Protocol Buffers (managed with buf), OpenAPI with ogen and chi
- PostgreSQL (pgx, goose migrations) and MongoDB
- zap for logging, caarlos0/env for configuration
- Docker, Docker Compose, Testcontainers
- GitHub Actions, golangci-lint, mockery, Task

## Getting started

**Requirements:** Go, Docker, [Task](https://taskfile.dev), `envsubst` (part of gettext).

```bash
# 1. Generate .env files from templates
task env:generate

# 2. Start the databases
task up-all

# 3. Run the services (from the repository root)
go run ./inventory/cmd
go run ./payment/cmd
go run ./order/cmd
```

Try it:

```bash
# Pick a part UUID
grpcurl -plaintext -d '{"filter":{}}' localhost:50051 inventory.v1.InventoryService/ListParts

# Create an order
curl -X POST localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"user_uuid":"<uuid>","part_uuids":["<part uuid>"]}'

# Pay for it
curl -X POST localhost:8080/api/v1/orders/<order uuid>/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"CARD"}'
```

An end-to-end smoke check of the whole chain is available as `task test-api`.

## Configuration

Every service reads its settings from environment variables. For local runs they are loaded from `deploy/compose/<service>/.env`, which `task env:generate` produces from `deploy/env/*.env.template`. In containers the same variables are passed through the environment. Required variables are validated at startup, so a missing value fails fast.

## Health checks and shutdown

- HTTP: `GET /health` on the order service.
- gRPC: the standard `grpc.health.v1.Health` service on inventory and payment.
- On `SIGINT` or `SIGTERM` each service stops accepting requests and then closes its resources in reverse order of creation (server, connections, databases, logger).

## Development

| Command                 | What it does                                 |
| ----------------------- | -------------------------------------------- |
| `task lint`             | Run golangci-lint for all modules            |
| `task test`             | Run unit tests                               |
| `task test-integration` | Run integration tests (needs Docker)         |
| `task test-coverage`    | Coverage for service and repository packages |
| `task mockery:gen`      | Regenerate mocks                             |
| `task gen`              | Regenerate protobuf and OpenAPI code         |

## Testing

- **Unit tests** cover services, repositories (with mocks) and converters.
- **Integration tests** use Testcontainers and the `integration` build tag, so `task test` never needs Docker:
  - order repository against a real PostgreSQL container with the real migrations;
  - inventory end-to-end test: the service is built from its Dockerfile, started next to a MongoDB container and called over gRPC like a black box.
- CI runs lint, unit tests and integration tests on every push; coverage is published to the badge above from `main`.

## Roadmap

- Asynchronous messaging with Kafka (order events between services)
- Redis cache for hot catalog reads

