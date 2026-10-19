package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger создает middleware для логирования времени выполнения запросов
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Засекаем время начала обработки запроса
		startTime := time.Now()

		// Логируем начало запроса
		log.Printf("⏱️ Начало запроса: %q %q", r.Method, r.URL.Path) //nolint:gosec // значения экранированы через %q

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Вычисляем время выполнения запроса
		duration := time.Since(startTime)

		// Логируем окончание запроса с временем выполнения
		log.Printf("✅ Запрос завершен: %q %q, время выполнения: %v", r.Method, r.URL.Path, duration) //nolint:gosec // значения экранированы через %q
	})
}
