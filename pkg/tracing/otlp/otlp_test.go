package otlp

import (
	"context"
	"testing"

	"github.com/efficientgo/core/testutil"

	"github.com/kubeservice-stack/common/pkg/tracing"

	"github.com/go-kit/log"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// This test creates an OTLP tracer, starts a span and checks whether it is logged in the exporter.
func TestContextTracing_ClientEnablesTracing(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()

	tracerOtel := newTraceProvider(context.Background(), tracesdk.NewSimpleSpanProcessor(exp), log.NewNopLogger(), "kubeservice", nil, tracesdk.AlwaysSample())
	tracer, _ := tracing.Bridge(tracerOtel, log.NewNopLogger())
	clientRoot, _ := tracing.StartSpan(tracing.ContextWithTracer(context.Background(), tracer), "a")

	testutil.Equals(t, 0, len(exp.GetSpans()))

	clientRoot.Finish()
	testutil.Equals(t, 1, len(exp.GetSpans()))
	testutil.Equals(t, 1, tracing.CountSampledSpans(exp.GetSpans()))
}
