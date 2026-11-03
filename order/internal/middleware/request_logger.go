package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/ivmerk/goshop/platform/pkg/logger"
)

// RequestLogger создает middleware для логирования времени выполнения запросов
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Засекаем время начала обработки запроса
		ctx := r.Context()
		startTime := time.Now()

		// Логируем начало запроса
		logger.Info(ctx, "request started", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Вычисляем время выполнения запроса
		duration := time.Since(startTime)

		// Логируем окончание запроса с временем выполнения
		logger.Info(ctx, "request finished", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("duration", duration.String()))
	})
}
