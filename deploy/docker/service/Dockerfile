# Универсальный Dockerfile для всех Go-сервисов монорепозитория.
# Сборка: docker build -f deploy/docker/service/Dockerfile --build-arg SERVICE=order .

# ============================
# Stage 1: Build
# ============================
FROM golang:1.27-alpine AS builder

# Имя собираемого сервиса (order, payment, assembly, inventory, iam, notification)
ARG SERVICE

RUN apk add --no-cache git

WORKDIR /app

# go.work и go.mod/go.sum всех модулей — слой с зависимостями кэшируется отдельно
COPY go.work go.work.sum ./
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/
COPY assembly/go.mod assembly/go.sum ./assembly/
COPY notification/go.mod notification/go.sum ./notification/
COPY iam/go.mod iam/go.sum ./iam/

RUN go mod download

# Исходники: сервис и общие модули
COPY platform ./platform
COPY shared ./shared
COPY ${SERVICE} ./${SERVICE}

# Собираем сервис (без CGO); миграции (если есть) складываем отдельно
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./${SERVICE}/cmd/main.go \
    && mkdir -p /out/migrations \
    && if [ -d "${SERVICE}/migrations" ]; then cp -r ${SERVICE}/migrations/. /out/migrations/; fi

# ============================
# Stage 2: Final image
# ============================
FROM alpine:3.21.3

ARG SERVICE

# Сертификаты нужны для HTTPS-клиентов (например, Telegram API)
RUN apk add --no-cache ca-certificates \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /out/app ./app
# Путь ./<service>/migrations совпадает со значением *_MIGRATION_DIRECTORY
COPY --from=builder /out/migrations ./${SERVICE}/migrations

USER appuser

ENTRYPOINT ["./app"]
