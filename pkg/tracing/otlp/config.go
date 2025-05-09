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
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

type retryConfig struct {
	RetryEnabled         bool          `yaml:"retry_enabled"`
	RetryInitialInterval time.Duration `yaml:"retry_initial_interval"`
	RetryMaxInterval     time.Duration `yaml:"retry_max_interval"`
	RetryMaxElapsedTime  time.Duration `yaml:"retry_max_elapsed_time"`
}

type Config struct {
	ClientType         string            `yaml:"client_type"`
	ServiceName        string            `yaml:"service_name"`
	ResourceAttributes map[string]string `yaml:"resource_attributes"`
	ReconnectionPeriod time.Duration     `yaml:"reconnection_period"`
	Compression        string            `yaml:"compression"`
	Insecure           bool              `yaml:"insecure"`
	Endpoint           string            `yaml:"endpoint"`
	URLPath            string            `yaml:"url_path"`
	Timeout            time.Duration     `yaml:"timeout"`
	RetryConfig        retryConfig       `yaml:"retry_config"`
	Headers            map[string]string `yaml:"headers"`
	TLSConfig          TLSConfig         `yaml:"tls_config"`
	SamplerType        string            `yaml:"sampler_type"`
	SamplerParam       string            `yaml:"sampler_param"`
}

func traceGRPCOptions(config Config) []otlptracegrpc.Option {
	var options []otlptracegrpc.Option
	if config.Endpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	if config.ReconnectionPeriod != 0 {
		options = append(options, otlptracegrpc.WithReconnectionPeriod(config.ReconnectionPeriod))
	}

	if config.Timeout != 0 {
		options = append(options, otlptracegrpc.WithTimeout(config.Timeout))
	}

	if config.Compression != "" {
		if config.Compression == "gzip" {
			options = append(options, otlptracegrpc.WithCompressor(config.Compression))
		}
	}

	if config.RetryConfig.RetryEnabled {
		options = append(options, otlptracegrpc.WithRetry(createGRPCRetryConfig(config)))
	}

	if config.Headers != nil {
		options = append(options, otlptracegrpc.WithHeaders(config.Headers))
	}

	return options
}

func traceHTTPOptions(config Config) []otlptracehttp.Option {
	var options []otlptracehttp.Option
	if config.Endpoint != "" {
		options = append(options, otlptracehttp.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	} else {
		tlsConfig, _ := NewTLSConfig(&config.TLSConfig)
		options = append(options, otlptracehttp.WithTLSClientConfig(tlsConfig))
	}

	if config.URLPath != "" {
		options = append(options, otlptracehttp.WithURLPath(config.URLPath))
	}

	if config.Compression != "" {
		if config.Compression == "gzip" {
			options = append(options, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
		}
	}

	if config.Timeout != 0 {
		options = append(options, otlptracehttp.WithTimeout(config.Timeout))
	}

	if config.RetryConfig.RetryEnabled {
		options = append(options, otlptracehttp.WithRetry(createHTTPRetryConfig(config)))
	}

	if config.Headers != nil {
		options = append(options, otlptracehttp.WithHeaders(config.Headers))
	}
	// how to specify JSON/binary format here?

	return options
}

func createHTTPRetryConfig(config Config) otlptracehttp.RetryConfig {

	var retryConfig otlptracehttp.RetryConfig
	if config.RetryConfig.RetryInitialInterval != 0 {
		retryConfig.InitialInterval = config.RetryConfig.RetryInitialInterval
	}

	if config.RetryConfig.RetryMaxInterval != 0 {
		retryConfig.MaxInterval = config.RetryConfig.RetryMaxInterval
	}

	if config.RetryConfig.RetryMaxElapsedTime != 0 {
		retryConfig.MaxElapsedTime = config.RetryConfig.RetryMaxElapsedTime
	}

	return retryConfig
}

func createGRPCRetryConfig(config Config) otlptracegrpc.RetryConfig {

	var retryConfig otlptracegrpc.RetryConfig
	if config.RetryConfig.RetryInitialInterval != 0 {
		retryConfig.InitialInterval = config.RetryConfig.RetryInitialInterval
	}

	if config.RetryConfig.RetryMaxInterval != 0 {
		retryConfig.MaxInterval = config.RetryConfig.RetryMaxInterval
	}

	if config.RetryConfig.RetryMaxElapsedTime != 0 {
		retryConfig.MaxElapsedTime = config.RetryConfig.RetryMaxElapsedTime
	}

	return retryConfig
}
