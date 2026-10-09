---

## 🛠 Актуальная структура проекта

```
.
├── README.md
├── Taskfile.yml
├── buf.work.yaml
├── go.work
├── go.work.sum
├── inventory
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── api
│       │   └── inventory
│       │       └── v1
│       │           ├── api.go
│       │           ├── get.go
│       │           └── list.go
│       ├── converter
│       │   └── part.go
│       ├── model
│       │   ├── errors.go
│       │   └── part.go
│       ├── repository
│       │   ├── converter
│       │   │   └── part.go
│       │   ├── mocks
│       │   │   └── mock_part_repository.go
│       │   ├── model
│       │   │   └── part.go
│       │   ├── part
│       │   │   ├── get.go
│       │   │   ├── init.go
│       │   │   ├── list.go
│       │   │   └── repository.go
│       │   └── repository.go
│       └── service
│           ├── mocks
│           │   └── mock_part_service.go
│           ├── part
│           │   ├── get.go
│           │   ├── get_test.go
│           │   ├── list.go
│           │   ├── list_test.go
│           │   ├── service.go
│           │   └── suite_test.go
│           └── service.go
├── order
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── api
│       │   └── order
│       │       └── v1
│       │           ├── api.go
│       │           ├── cancel.go
│       │           ├── create.go
│       │           ├── get.go
│       │           ├── new_order.go
│       │           └── pay.go
│       ├── client
│       │   ├── converter
│       │   │   └── part.go
│       │   └── grpc
│       │       ├── client.go
│       │       ├── inventory
│       │       │   └── v1
│       │       │       ├── client.go
│       │       │       └── list_parts.go
│       │       ├── mocks
│       │       │   ├── mock_inventory_client.go
│       │       │   └── mock_payment_client.go
│       │       └── payment
│       │           └── v1
│       │               ├── client.go
│       │               └── pay_order.go
│       ├── converter
│       │   └── order.go
│       ├── model
│       │   ├── error.go
│       │   ├── order.go
│       │   └── part.go
│       ├── repository
│       │   ├── converter
│       │   │   └── order.go
│       │   ├── mocks
│       │   │   └── mock_order_repository.go
│       │   ├── model
│       │   │   └── order.go
│       │   ├── order
│       │   │   ├── create.go
│       │   │   ├── get.go
│       │   │   ├── repository.go
│       │   │   └── update.go
│       │   └── repository.go
│       └── service
│           ├── mocks
│           │   └── mock_order_service.go
│           ├── order
│           │   ├── cancel.go
│           │   ├── cancel_test.go
│           │   ├── create.go
│           │   ├── create_test.go
│           │   ├── get.go
│           │   ├── get_test.go
│           │   ├── pay.go
│           │   ├── pay_test.go
│           │   ├── service.go
│           │   └── suite_test.go
│           └── service.go
├── package-lock.json
├── package.json
├── payment
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── api
│       │   └── payment
│       │       └── v1
│       │           ├── api.go
│       │           └── pay.go
│       ├── model
│       │   └── errors.go
│       └── service
│           ├── mocks
│           │   └── mock_payment_service.go
│           ├── payment
│           │   ├── pay.go
│           │   ├── pay_test.go
│           │   ├── service.go
│           │   └── suite_test.go
│           └── service.go
└── shared
    ├── api
    │   └── order
    │       └── v1
    │           ├── components/ (create_order_request, create_order_response, enums/, errors/, get_order_response, order_dto, pay_order_request, pay_order_response)
    │           ├── order.openapi.yaml
    │           ├── params/order_uuid.yaml
    │           └── paths/ (order_by_uuid, order_cancel, order_pay, orders)
    ├── go.mod
    ├── go.sum
    ├── pkg
    │   ├── openapi/order/v1/   — ogen-генерация (oas_*_gen.go)
    │   └── proto
    │       ├── inventory/v1/   — inventory.pb.go, inventory_grpc.pb.go
    │       └── payment/v1/     — payment.pb.go, payment_grpc.pb.go
    └── proto
        ├── buf.gen.yaml
        ├── buf.yaml
        ├── inventory/v1/inventory.proto
        └── payment/v1/payment.proto
```

---

### 🔧 Комментарии

- **Каждый сервис — отдельный модуль с `go.mod`**, подключённый в `go.work`.
- **В каждом сервисе выделены слои: `api`, `service`, `repository`.**
- **Контракты и автогенерация (`protobuf`, `ogen`) находятся в `shared/`**.
- **Моки и тесты** для всех слоёв размещаются рядом с реализациями.

---

## 💡 Полезные подсказки

- 📚 Прежде чем писать код, **ознакомьтесь с уроками второй недели** — они показывают, как правильно выделить слои и писать тесты.
- 🧠 Если вы застряли — **обратитесь в чат или к ревьюеру**, особенно если за 30 минут не продвинулись.
> **Спросить — не значит сдаться.** Это ускоряет обучение и избавляет от тупиков.

---

**Автор курса: Олег Козырев, 2025**
