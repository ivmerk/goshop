package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	httpRequests = promauto.NewCounterVec(prom.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество обработанных HTTP-запросов.",
	}, []string{"method", "status"})

	httpDuration = promauto.NewHistogramVec(prom.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Длительность обработки HTTP-запросов.",
		Buckets: prom.DefBuckets,
	}, []string{"method"})

	grpcRequests = promauto.NewCounterVec(prom.CounterOpts{
		Name: "grpc_server_requests_total",
		Help: "Общее количество обработанных unary gRPC-вызовов.",
	}, []string{"method", "code"})

	grpcDuration = promauto.NewHistogramVec(prom.HistogramOpts{
		Name:    "grpc_server_request_duration_seconds",
		Help:    "Длительность обработки unary gRPC-вызовов.",
		Buckets: prom.DefBuckets,
	}, []string{"method"})
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// HTTPMiddleware считает запросы и их длительность. Путь в метки не входит
// сознательно: с path-параметрами (uuid) он взорвал бы кардинальность метрик.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		httpRequests.WithLabelValues(r.Method, strconv.Itoa(rec.status)).Inc()
		httpDuration.WithLabelValues(r.Method).Observe(time.Since(start).Seconds())
	})
}

// UnaryServerInterceptor считает unary gRPC-вызовы по методу и коду ответа.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		grpcRequests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
		grpcDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
