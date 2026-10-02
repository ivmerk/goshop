package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ivmerk/goshop/order/internal/models"

	customMiddleware "github.com/ivmerk/goshop/order/internal/middleware"
	orderV1 "github.com/ivmerk/goshop/shared/pkg/openapi/order/v1"
)

const (
	httpPort     = "8080"
	urlParamCity = "city"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderHandler struct {
	storage   *models.OrderStorage
	payment   PaymentClient
	inventory InventoryClient
}

type InventoryClient interface {
	ListPartPrices(ctx context.Context, uuids []string) (map[string]float64, error)
}

func NewOrderHandler(storage *models.OrderStorage, inventory InventoryClient, payment PaymentClient) *OrderHandler {
	return &OrderHandler{storage: storage, inventory: inventory, payment: payment}
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	switch order.Status {
	case models.OrderStatusPaid:
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "order is already paid and cannot be canceled",
		}, nil
	case models.OrderStatusPendingPayment:
		order.Status = models.OrderStatusCancelled
		h.storage.UpdateOrder(order)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, method models.PaymentMethod) (transactionUUID string, err error)
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{Code: http.StatusNotFound, Message: "order not found"}, nil
	}

	if req.PaymentMethod == orderV1.PaymentMethodUNKNOWN {
		return &orderV1.BadRequestError{Code: http.StatusBadRequest, Message: "unknown payment method"}, nil
	}

	method := models.PaymentMethod(req.PaymentMethod)
	transactionUUID, err := h.payment.PayOrder(ctx, order.UUID, order.User, method)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "payment failed"}, nil
	}

	order.Status = models.OrderStatusPaid
	order.Transaction = transactionUUID
	order.Payment = method
	h.storage.UpdateOrder(order)

	txID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction id"}, nil
	}
	return &orderV1.PayOrderResponse{TransactionUUID: txID}, nil
}

func (h *OrderHandler) GetOrderById(_ context.Context, params orderV1.GetOrderByIdParams) (orderV1.GetOrderByIdRes, error) {
	order, ok := h.storage.GetOrder(params.OrderUUID.String())
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	userUUID, err := uuid.Parse(order.User)
	if err != nil {
		return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid user uuid"}, nil
	}

	partUUIDs := make([]uuid.UUID, 0, len(order.Parts))
	for _, p := range order.Parts {
		id, err := uuid.Parse(p)
		if err != nil {
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid part uuid"}, nil
		}
		partUUIDs = append(partUUIDs, id)
	}

	resp := &orderV1.GetOrderResponse{
		OrderUUID:  params.OrderUUID,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: order.Total,
		Status:     orderV1.OrderStatus(order.Status),
	}

	// Эти поля есть только у оплаченного заказа.
	if order.Transaction != "" {
		txID, err := uuid.Parse(order.Transaction)
		if err != nil {
			return &orderV1.InternalServerError{Code: http.StatusInternalServerError, Message: "invalid transaction uuid"}, nil
		}
		resp.TransactionUUID = orderV1.NewOptUUID(txID)
	}
	if order.Payment != "" {
		resp.PaymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(order.Payment))
	}

	return resp, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "part_uuids must not be empty",
		}, nil
	}

	partUUIDs := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		partUUIDs = append(partUUIDs, id.String())
	}

	prices, err := h.inventory.ListPartPrices(ctx, partUUIDs)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get parts from inventory",
		}, nil
	}

	var total float64
	for _, id := range partUUIDs {
		price, ok := prices[id]
		if !ok {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "part not found: " + id,
			}, nil
		}
		total += price
	}

	orderUUID := uuid.New()
	h.storage.AddOrder(&models.Order{
		UUID:   orderUUID.String(),
		User:   req.UserUUID.String(),
		Parts:  partUUIDs,
		Total:  total,
		Status: models.OrderStatusPendingPayment,
	})

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total,
	}, nil
}

func main() {
	storage := models.NewOrderStorage()

	orderHandler := NewOrderHandler(storage, nil, nil)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
