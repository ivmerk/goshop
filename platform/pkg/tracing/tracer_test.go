package tracing

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestInit_DisabledWithoutEndpoint(t *testing.T) {
	t.Setenv(EnvEndpoint, "")

	require.NoError(t, Init(context.Background(), "svc"))
}

func TestPropagation_TraceContextCrossesServices(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	tp := newProvider("order", sdktrace.WithSyncer(exp))
	setGlobal(tp)

	// «Сервис A» создаёт спан и кладёт trace context в заголовки исходящего запроса.
	ctx, parent := Tracer("test").Start(context.Background(), "parent")
	headers := http.Header{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
	parent.End()
	require.NotEmpty(t, headers.Get("traceparent"))

	// «Сервис B» извлекает контекст из заголовков: его спан должен остаться в том же трейсе.
	remote := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(headers))
	_, child := Tracer("test").Start(remote, "child")
	child.End()

	spans := exp.GetSpans()
	require.Len(t, spans, 2)
	require.Equal(t, spans[0].SpanContext.TraceID(), spans[1].SpanContext.TraceID())
	require.Equal(t, spans[0].SpanContext.SpanID(), spans[1].Parent.SpanID())
}
