/*
Copyright 2025 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jaeger

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/kubeservice-stack/common/pkg/tracing"
)

var parentConfig = ParentBasedSamplerConfig{LocalParentSampled: true}

// This test shows that if sample factor will enable tracing on client process, even when it would be disabled on server
// it will be still enabled for all spans within this span.
func TestContextTracing_ClientEnablesProbabilisticTracing(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	config := Config{
		SamplerType:         "probabilistic",
		SamplerParam:        1.0,
		SamplerParentConfig: parentConfig,
	}
	sampler := getSampler(config)

	tracerOtel := newTraceProvider(
		context.Background(),
		"tracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler,
		[]attribute.KeyValue{},
	)
	tracer, _ := tracing.Bridge(tracerOtel, log.NewNopLogger())
	clientRoot, clientCtx := tracing.StartSpan(tracing.ContextWithTracer(context.Background(), tracer), "a")

	config.SamplerParam = 0.0
	sampler2 := getSampler(config)
	// Simulate Server process with different tracer, but with client span in context.
	srvTracerOtel := newTraceProvider(
		context.Background(),
		"srvTracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler2, // never sample
		[]attribute.KeyValue{},
	)
	srvTracer, _ := tracing.Bridge(srvTracerOtel, log.NewNopLogger())

	srvRoot, srvCtx := tracing.StartSpan(tracing.ContextWithTracer(clientCtx, srvTracer), "b")
	srvChild, _ := tracing.StartSpan(srvCtx, "bb")

	tracing.CountSpans_ClientEnablesTracing(t, exp, clientRoot, srvRoot, srvChild)
}

// This test shows that if sample factor will disable tracing on client process,  when it would be enabled on server
// it will be still disabled for all spans within this span.
func TestContextTracing_ClientDisablesProbabilisticTracing(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()

	config := Config{
		SamplerType:         "probabilistic",
		SamplerParam:        0.0,
		SamplerParentConfig: parentConfig,
	}
	sampler := getSampler(config)
	tracerOtel := newTraceProvider(
		context.Background(),
		"tracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler, // never sample
		[]attribute.KeyValue{},
	)
	tracer, _ := tracing.Bridge(tracerOtel, log.NewNopLogger())

	clientRoot, clientCtx := tracing.StartSpan(tracing.ContextWithTracer(context.Background(), tracer), "a")

	config.SamplerParam = 1.0
	sampler2 := getSampler(config)
	// Simulate Server process with different tracer, but with client span in context.
	srvTracerOtel := newTraceProvider(
		context.Background(),
		"srvTracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler2, // never sample
		[]attribute.KeyValue{},
	)
	srvTracer, _ := tracing.Bridge(srvTracerOtel, log.NewNopLogger())

	srvRoot, srvCtx := tracing.StartSpan(tracing.ContextWithTracer(clientCtx, srvTracer), "b")
	srvChild, _ := tracing.StartSpan(srvCtx, "bb")

	tracing.ContextTracing_ClientDisablesTracing(t, exp, clientRoot, srvRoot, srvChild)
}

func TestContextTracing_ClientDisablesAlwaysOnSampling(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()

	config := Config{
		SamplerType:  SamplerTypeConstant,
		SamplerParam: 0,
	}
	sampler := getSampler(config)
	tracerOtel := newTraceProvider(
		context.Background(),
		"tracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler, // never sample
		[]attribute.KeyValue{},
	)
	tracer, _ := tracing.Bridge(tracerOtel, log.NewNopLogger())

	clientRoot, clientCtx := tracing.StartSpan(tracing.ContextWithTracer(context.Background(), tracer), "a")

	config.SamplerParam = 1
	sampler2 := getSampler(config)
	// Simulate Server process with different tracer, but with client span in context.
	srvTracerOtel := newTraceProvider(
		context.Background(),
		"srvTracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler2, // never sample
		[]attribute.KeyValue{},
	)
	srvTracer, _ := tracing.Bridge(srvTracerOtel, log.NewNopLogger())

	srvRoot, srvCtx := tracing.StartSpan(tracing.ContextWithTracer(clientCtx, srvTracer), "b")
	srvChild, _ := tracing.StartSpan(srvCtx, "bb")

	tracing.ContextTracing_ClientDisablesTracing(t, exp, clientRoot, srvRoot, srvChild)
}

// This test shows that if span will contain special baggage (for example from special HTTP header), even when sample
// factor will disable client & server tracing, it will be still enabled for all spans within this span.
func TestContextTracing_ForceTracing(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	config := Config{
		SamplerType:         "probabilistic",
		SamplerParam:        0.0,
		SamplerParentConfig: parentConfig,
	}
	sampler := getSampler(config)
	tracerOtel := newTraceProvider(
		context.Background(),
		"tracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler,
		[]attribute.KeyValue{},
	)
	tracer, _ := tracing.Bridge(tracerOtel, log.NewNopLogger())

	// Start the root span with the tag to force tracing.
	clientRoot, clientCtx := tracing.StartSpan(
		tracing.ContextWithTracer(context.Background(), tracer),
		"a",
		opentracing.Tag{Key: tracing.ForceTracingAttributeKey, Value: "true"},
	)

	// Simulate Server process with different tracer, but with client span in context.
	srvTracerOtel := newTraceProvider(
		context.Background(),
		"srvTracerOtel",
		log.NewNopLogger(),
		tracesdk.NewSimpleSpanProcessor(exp),
		sampler,
		[]attribute.KeyValue{},
	)
	srvTracer, _ := tracing.Bridge(srvTracerOtel, log.NewNopLogger())

	srvRoot, srvCtx := tracing.StartSpan(tracing.ContextWithTracer(clientCtx, srvTracer), "b")
	srvChild, _ := tracing.StartSpan(srvCtx, "bb")

	tracing.ContextTracing_ForceTracing(t, exp, clientRoot, srvRoot, srvChild)
}

func TestParseTags(t *testing.T) {
	for _, tcase := range []struct {
		input    string
		expected []attribute.KeyValue
	}{
		{
			input:    "key=value",
			expected: []attribute.KeyValue{attribute.String("key", "value")},
		},
		{
			input: "key1=value1,key2=value2",
			expected: []attribute.KeyValue{attribute.String("key1", "value1"),
				attribute.String("key2", "value2")},
		},
		{
			input:    "",
			expected: []attribute.KeyValue{},
		},
		{
			// Incorrectly formatted string with leading comma still yields the right tags.
			input:    ",key=value",
			expected: []attribute.KeyValue{attribute.String("key", "value")},
		},
		{
			// Incorrectly formatted string with trailing comma still yields the right tags.
			input:    "key=value,",
			expected: []attribute.KeyValue{attribute.String("key", "value")},
		},
		{
			// Leading and trailing spaces in tags are trimmed.
			input:    " key=value  ",
			expected: []attribute.KeyValue{attribute.String("key", "value")},
		},
		{
			input:    "key=${env:default_val}",
			expected: []attribute.KeyValue{attribute.String("key", "default_val")},
		},
	} {
		if ok := t.Run("", func(t *testing.T) {
			exists := false
			envVal := ""
			envVar := ""
			// Check if env vars are used.
			if strings.Contains(tcase.input, "${") {
				envVal, envVar, exists = extractValueOfEnvVar(tcase.input)
				// Set a temporary value just for testing.
				tempEnvVal := "temp_val"
				os.Setenv(envVar, tempEnvVal)
				tcase.expected = []attribute.KeyValue{attribute.String("key", tempEnvVal)}
			}
			attrs := parseTags(tcase.input)
			testutil.Equals(t, tcase.expected, attrs)

			// Reset the env var to the old value, if needed.
			if exists {
				os.Setenv(envVar, envVal)
			}
		}); !ok {
			return
		}
	}
}

func extractValueOfEnvVar(input string) (string, string, bool) {
	kv := strings.SplitN(input, "=", 2)
	_, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

	if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
		ed := strings.SplitN(v[2:len(v)-1], ":", 2)
		e, d := ed[0], ed[1]
		envVal, exists := os.LookupEnv(e)
		if !exists {
			return d, e, exists
		}
		return envVal, e, exists
	}

	return "", "", false
}
