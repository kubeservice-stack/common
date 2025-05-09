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

package otlp

import (
	"context"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/attribute"

	"github.com/kubeservice-stack/common/pkg/tracing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	_ "google.golang.org/grpc/encoding/gzip"
	"gopkg.in/yaml.v2"
)

const (
	TracingClientGRPC                  string = "grpc"
	TracingClientHTTP                  string = "http"
	AlwaysSample                       string = "alwayssample"
	NeverSample                        string = "neversample"
	TraceIDRatioBasedSample            string = "traceidratiobased"
	ParentBasedAlwaysSample            string = "parentbasedalwayssample"
	ParentBasedNeverSample             string = "parentbasedneversample"
	ParentBasedTraceIDRatioBasedSample string = "parentbasedtraceidratiobased"
)

// NewOTELTracer returns an OTLP exporter based tracer.
func NewTracerProvider(ctx context.Context, logger log.Logger, conf []byte) (*tracesdk.TracerProvider, error) {
	config := Config{}
	if err := yaml.Unmarshal(conf, &config); err != nil {
		return nil, err
	}

	var exporter *otlptrace.Exporter
	var err error
	switch strings.ToLower(config.ClientType) {
	case TracingClientHTTP:
		options := traceHTTPOptions(config)

		client := otlptracehttp.NewClient(options...)
		exporter, err = otlptrace.New(ctx, client)
		if err != nil {
			return nil, err
		}

	case TracingClientGRPC:
		options := traceGRPCOptions(config)
		client := otlptracegrpc.NewClient(options...)
		exporter, err = otlptrace.New(ctx, client)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("otlp: invalid client type. Only 'http' and 'grpc' are accepted. ")
	}

	processor := tracesdk.NewBatchSpanProcessor(exporter)
	sampler, err := getSampler(config)
	if err != nil {
		logger.Log(err)
	}
	tp := newTraceProvider(ctx, processor, logger, config.ServiceName, config.ResourceAttributes, sampler)

	return tp, nil
}

func newTraceProvider(
	ctx context.Context,
	processor tracesdk.SpanProcessor,
	logger log.Logger,
	serviceName string,
	attrs map[string]string,
	sampler tracesdk.Sampler,
) *tracesdk.TracerProvider {
	resourceAttrs := make([]attribute.KeyValue, 0, len(attrs)+1)
	if serviceName != "" {
		resourceAttrs = append(resourceAttrs, semconv.ServiceNameKey.String(serviceName))
	}
	for k, v := range attrs {
		resourceAttrs = append(resourceAttrs, attribute.String(k, v))
	}
	r, err := resource.New(ctx, resource.WithAttributes(resourceAttrs...))
	if err != nil {
		level.Warn(logger).Log("msg", "jaeger: detecting resources for tracing provider failed", "err", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(processor),
		tracesdk.WithResource(r),
		tracesdk.WithSampler(
			tracing.SamplerWithOverride(
				sampler, tracing.ForceTracingAttributeKey,
			),
		),
	)
	return tp
}

func getSampler(config Config) (tracesdk.Sampler, error) {
	switch strings.ToLower(config.SamplerType) {
	case AlwaysSample:
		return tracesdk.AlwaysSample(), nil
	case NeverSample:
		return tracesdk.NeverSample(), nil
	case TraceIDRatioBasedSample:
		arg, err := strconv.ParseFloat(config.SamplerParam, 64)
		if err != nil {
			return tracesdk.TraceIDRatioBased(1.0), err
		}
		return tracesdk.TraceIDRatioBased(arg), nil
	case ParentBasedAlwaysSample:
		return tracesdk.ParentBased(tracesdk.AlwaysSample()), nil
	case ParentBasedNeverSample:
		return tracesdk.ParentBased(tracesdk.NeverSample()), nil
	case ParentBasedTraceIDRatioBasedSample:
		arg, err := strconv.ParseFloat(config.SamplerParam, 64)
		if err != nil {
			return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0)), err
		}
		return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(arg)), nil
	}

	return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0)), nil
}
