// Package tracing инициализирует OpenTelemetry: провайдер трейсов, OTLP-экспортёр
// и propagator (W3C Trace Context), чтобы trace context передавался между сервисами.
package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ivmerk/goshop/platform/pkg/closer"
)

const (
	// EnvEndpoint — адрес OTLP/gRPC-приёмника (стандартная переменная OpenTelemetry),
	// например "http://jaeger:4317". Схема http:// означает соединение без TLS.
	// Если переменная пуста, трейсинг отключён: сервис работает как раньше.
	EnvEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"

	shutdownTimeout = 5 * time.Second
	batchTimeout    = time.Second
)

// Init настраивает глобальный TracerProvider для сервиса serviceName и регистрирует
// его остановку в closer (при остановке сбрасывает неотправленные спаны).
// Без OTEL_EXPORTER_OTLP_ENDPOINT ничего не делает и возвращает nil.
func Init(ctx context.Context, serviceName string) error {
	endpoint := os.Getenv(EnvEndpoint)
	if endpoint == "" {
		return nil
	}

	// Экспортёр сам читает OTEL_EXPORTER_OTLP_ENDPOINT, повторно адрес не передаём.
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return fmt.Errorf("create otlp trace exporter: %w", err)
	}

	tp := newProvider(serviceName, sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(batchTimeout)))
	setGlobal(tp)

	closer.AddNamed("tracer provider", func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		return tp.Shutdown(ctx)
	})

	return nil
}

// newProvider создаёт провайдер с ресурсом сервиса; опции позволяют подставить экспортёр (в тестах).
func newProvider(serviceName string, opts ...sdktrace.TracerProviderOption) *sdktrace.TracerProvider {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		attribute.String("deployment.environment", "local"),
	)

	return sdktrace.NewTracerProvider(append(opts, sdktrace.WithResource(res))...)
}

// setGlobal делает провайдер глобальным и включает W3C Trace Context + Baggage:
// именно по ним otelhttp/otelgrpc передают и принимают trace context между сервисами.
func setGlobal(tp *sdktrace.TracerProvider) {
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

// Tracer возвращает именованный tracer для ручных спанов.
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
