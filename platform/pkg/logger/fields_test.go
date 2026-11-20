package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func fieldMap(fields []zap.Field) map[string]string {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f.Key] = f.String
	}

	return m
}

func TestFieldsFromContext_TakesIDsFromActiveSpan(t *testing.T) {
	traceID, err := trace.TraceIDFromHex("5b8efff798038103d269b633813fc60c")
	require.NoError(t, err)
	spanID, err := trace.SpanIDFromHex("eee19b7ec3c1b174")
	require.NoError(t, err)

	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled,
	}))

	got := fieldMap(fieldsFromContext(ctx))

	require.Equal(t, "5b8efff798038103d269b633813fc60c", got["trace_id"])
	require.Equal(t, "eee19b7ec3c1b174", got["span_id"])
}

func TestFieldsFromContext_FallsBackToManualTraceID(t *testing.T) {
	ctx := context.WithValue(context.Background(), traceIDKey, "manual-id")

	got := fieldMap(fieldsFromContext(ctx))

	require.Equal(t, "manual-id", got["trace_id"])
	require.NotContains(t, got, "span_id")
}

func TestFieldsFromContext_EmptyWithoutTrace(t *testing.T) {
	require.Empty(t, fieldsFromContext(context.Background()))
}
