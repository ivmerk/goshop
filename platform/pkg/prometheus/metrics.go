// Package prometheus — общая обвязка метрик для всех сервисов:
// отдельный HTTP-сервер с /metrics и middleware для HTTP и gRPC.
package prometheus

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ivmerk/goshop/platform/pkg/closer"
)

const (
	// EnvAddr — переменная окружения с адресом сервера метрик.
	EnvAddr = "METRICS_ADDR"
	// DefaultAddr — адрес по умолчанию.
	DefaultAddr = ":9100"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// Serve запускает в фоне HTTP-сервер с /metrics (глобальный реестр, включая
// метрики Go runtime и процесса) и регистрирует его остановку в closer.
// Сервер намеренно отдельный от основного: /metrics не проходит авторизацию по сессии.
func Serve() error {
	addr := os.Getenv(EnvAddr)
	if addr == "" {
		addr = DefaultAddr
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Метрики вторичны: падение их сервера не должно ронять сервис.
			_ = err
		}
	}()

	closer.AddNamed("metrics server", func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	return nil
}
