![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ivmerk/46bf0f7a10ffaf0e26dc0646ad729a46/raw/coverage.json)

# goshop

A small spare-parts ordering platform built as a set of Go microservices. Customers sign in, create an order for a list of parts, the system prices it from the parts catalog, and the order is then paid through a payment service, assembled asynchronously and the customer is notified in Telegram.

## Services

| Service       | Transport      | Storage    | Responsibility                                           |
| ------------- | -------------- | ---------- | -------------------------------------------------------- |
| **iam**       | gRPC           | PostgreSQL, Redis | Users and sessions: register, login, whoami, get user |
| **inventory** | gRPC           | MongoDB    | Parts catalog: get a part, list parts with filters       |
| **order**     | REST (OpenAPI) | PostgreSQL | Create, read and pay orders; calls inventory and payment |
| **payment**   | gRPC           | none       | Processes a payment and returns a transaction ID         |
| **assembly**  | Kafka          | none       | Consumes `OrderPaid`, simulates assembly (1–10 s), produces `ShipAssembled` |
| **notification** | Kafka, Telegram | none  | Consumes `OrderPaid` and `ShipAssembled`, sends Telegram messages |

```
client ──REST──▶ order ──gRPC──▶ inventory ──▶ MongoDB
                  │                  │
                  ├──gRPC──▶ payment │
                  ├──▶ PostgreSQL    │
                  └──Whoami──▶ iam ◀─┘ Whoami
                               ├──▶ PostgreSQL (users)
                               └──▶ Redis (sessions)

order ──OrderPaid (order.paid)──▶ Kafka ──▶ assembly ──ShipAssembled (order.assembled)──▶ Kafka ──▶ order
                                   └──────────────▶ notification ◀──────────────┘  └──▶ Telegram
```

Order flow: the client registers and logs in through iam and gets a session UUID. `POST /api/v1/orders` takes the user from the session (not from the body), validates the parts through inventory and stores the order as `PENDING_PAYMENT`. `POST /api/v1/orders/{uuid}/pay` calls payment and moves the order to `PAID` and publishes `OrderPaid`. Assembly picks it up and publishes `ShipAssembled`; order consumes it and moves the order to `ASSEMBLED` (such an order can no longer be cancelled). Notification sends a Telegram message for both events.

### Authentication

iam exposes two gRPC services: `user.v1.UserService` (`Register`, `GetUser`) and `auth.v1.AuthService` (`Login`, `Whoami`). Passwords are stored as bcrypt hashes in PostgreSQL (migrations are applied on startup). `Login` creates a session in Redis under `session:<uuid>` with a 24 h TTL (`SESSION_TTL`) and returns its UUID.

The session UUID travels with every request and each service validates it through `AuthService.Whoami`:

| Hop                          | Where the session lives                 | Check                                      |
| ---------------------------- | --------------------------------------- | ------------------------------------------ |
| client → order (REST)        | header `X-Session-Uuid`                 | HTTP middleware, 401 `MISSING_SESSION` / `INVALID_SESSION`; `/health` is open |
| order → inventory (gRPC)     | metadata `session-uuid`                 | order forwards it with `ForwardSessionUUIDToGRPC` |
| client → inventory (gRPC)    | metadata `session-uuid`                 | unary interceptor, `Unauthenticated` on failure |

The middleware and the interceptor live in `platform/pkg/middleware`. They put the user (`UserInfo`) and the session UUID into the request context. Unknown login and wrong password are indistinguishable (`Unauthenticated`), so logins cannot be enumerated.

Events are protobuf messages keyed by order UUID, delivered at-least-once (consumers are idempotent; a transactional outbox is not implemented).

## Repository layout

```
.
├── iam/            # gRPC service: users (PostgreSQL), sessions (Redis)
├── inventory/      # gRPC service, MongoDB
├── order/          # REST service, PostgreSQL, gRPC clients
├── payment/        # gRPC service
├── assembly/       # Kafka consumer/producer, simulates order assembly
├── notification/   # Kafka consumer, Telegram bot
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
- PostgreSQL (pgx, goose migrations), MongoDB and Redis (redigo)
- bcrypt for password hashing
- zap for logging, caarlos0/env for configuration
- Kafka (KRaft, sarama) with protobuf events
- Telegram Bot API (go-telegram/bot)
- Docker, Docker Compose, Testcontainers
- GitHub Actions, golangci-lint, mockery, Task

## Getting started

**Requirements:** Go, Docker, [Task](https://taskfile.dev), `envsubst` (part of gettext).

```bash
# 1. Generate .env files from templates
task env:generate

# 2. Start the databases (iam: PostgreSQL + Redis)
task up-all

# 3. Run the services (from the repository root)
go run ./iam/cmd
go run ./inventory/cmd
go run ./payment/cmd
go run ./order/cmd
go run ./assembly/cmd
go run ./notification/cmd   # needs a Telegram bot token, see below
```

Try it:

```bash
# Register and log in (the response contains sessionUuid)
grpcurl -plaintext -d '{"login":"john","password":"SecurePassword123!","email":"john@example.com"}' \
  localhost:50053 user.v1.UserService/Register
grpcurl -plaintext -d '{"login":"john","password":"SecurePassword123!"}' \
  localhost:50053 auth.v1.AuthService/Login

# Pick a part UUID
grpcurl -plaintext -H "session-uuid: <session uuid>" -d '{"filter":{}}' \
  localhost:50051 inventory.v1.InventoryService/ListParts

# Create an order
curl -X POST localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "X-Session-Uuid: <session uuid>" \
  -d '{"part_uuids":["<part uuid>"]}'

# Pay for it
curl -X POST localhost:8080/api/v1/orders/<order uuid>/pay \
  -H "Content-Type: application/json" \
  -H "X-Session-Uuid: <session uuid>" \
  -d '{"payment_method":"CARD"}'
```

An end-to-end smoke check of the whole chain is available as `task test-api` (needs all services running). It runs `task test-iam` (register, login, whoami), then `task test-api-auth-rejection` (requests without or with an unknown session are rejected by inventory and order) and finally the full order flow with a valid session.

## Configuration

Every service reads its settings from environment variables. For local runs they are loaded from `deploy/compose/<service>/.env`, which `task env:generate` produces from `deploy/env/*.env.template`. In containers the same variables are passed through the environment. Required variables are validated at startup, so a missing value fails fast.

### Telegram bot

Create a bot via @BotFather, write `/start` to it and take your chat id from `https://api.telegram.org/bot<TOKEN>/getUpdates`. Put them into `deploy/env/.env` (git-ignored) as `NOTIFICATION_TELEGRAM_BOT_TOKEN` and `NOTIFICATION_TELEGRAM_CHAT_ID`, then run `task env:generate`. Never commit the token.

## Health checks and shutdown

- HTTP: `GET /health` on the order service.
- gRPC: the standard `grpc.health.v1.Health` service on iam, inventory and payment. Inventory requires a session for every unary call, including the health check.
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

- **Unit tests** cover services, repositories (with mocks) and converters; iam services are tested with mocked repositories.
- **Integration tests** use Testcontainers and the `integration` build tag, so `task test` never needs Docker:
  - order repository against a real PostgreSQL container with the real migrations;
  - inventory end-to-end test: the service is built from its Dockerfile, started next to a MongoDB container and called over gRPC like a black box. IAM is replaced by an in-process stub `AuthService` on the host (reachable from the container as `host.docker.internal`), and the suite also checks that calls without a valid session are rejected.
- CI runs lint, unit tests and integration tests on every push; coverage is published to the badge above from `main`.

## Roadmap

- Redis cache for hot catalog reads

