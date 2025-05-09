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

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel/attribute"
	otel_jaeger "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"gopkg.in/yaml.v2"

	"github.com/kubeservice-stack/common/pkg/tracing"
)

// NewTracerProvider returns a new instance of an OpenTelemetry tracer provider.
func NewTracerProvider(ctx context.Context, logger log.Logger, conf []byte) (*tracesdk.TracerProvider, error) {
	config := Config{}
	if err := yaml.Unmarshal(conf, &config); err != nil {
		return nil, err
	}

	printDeprecationWarnings(config, logger)

	var exporter *otel_jaeger.Exporter
	var err error

	if config.Endpoint != "" {
		collectorOptions := getCollectorEndpoints(config)

		exporter, err = otel_jaeger.New(otel_jaeger.WithCollectorEndpoint(collectorOptions...))
		if err != nil {
			return nil, err
		}
	} else if config.AgentHost != "" && config.AgentPort != 0 {
		jaegerAgentEndpointOptions := getAgentEndpointOptions(config)

		exporter, err = otel_jaeger.New(otel_jaeger.WithAgentEndpoint(jaegerAgentEndpointOptions...))
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err = otel_jaeger.New(otel_jaeger.WithAgentEndpoint())
		if err != nil {
			return nil, err
		}
	}

	var tags []attribute.KeyValue
	if config.Tags != "" {
		tags = getAttributesFromTags(config)
	}

	sampler := getSampler(config)
	var processorOptions []tracesdk.BatchSpanProcessorOption
	var processor tracesdk.SpanProcessor
	if config.ReporterMaxQueueSize != 0 {
		processorOptions = append(processorOptions, tracesdk.WithMaxQueueSize(config.ReporterMaxQueueSize))
	}

	//Ref: https://epsagon.com/observability/opentelemetry-best-practices-overview-part-2-2/ .
	if config.ReporterFlushInterval != 0 {
		processorOptions = append(processorOptions, tracesdk.WithBatchTimeout(config.ReporterFlushInterval))
	}

	processor = tracesdk.NewBatchSpanProcessor(exporter, processorOptions...)

	tp := newTraceProvider(ctx, config.ServiceName, logger, processor, sampler, tags)

	return tp, nil
}

// getAttributesFromTags returns tags as OTel attributes.
func getAttributesFromTags(config Config) []attribute.KeyValue {
	return parseTags(config.Tags)
}

func newTraceProvider(ctx context.Context, serviceName string, logger log.Logger, processor tracesdk.SpanProcessor,
	sampler tracesdk.Sampler, tags []attribute.KeyValue) *tracesdk.TracerProvider {

	resource, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
		resource.WithAttributes(tags...),
	)
	if err != nil {
		level.Warn(logger).Log("msg", "jaeger: detecting resources for tracing provider failed", "err", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSpanProcessor(processor),
		tracesdk.WithSampler(
			tracing.SamplerWithOverride(
				sampler, tracing.ForceTracingAttributeKey,
			),
		),
		tracesdk.WithResource(resource),
	)

	return tp
}
