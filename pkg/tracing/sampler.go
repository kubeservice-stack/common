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

package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// ForceTracingAttributeKey is used to signalize a span should be traced.
const ForceTracingAttributeKey = "force_tracing"

type samplerWithOverride struct {
	baseSampler tracesdk.Sampler
	overrideKey attribute.Key
}

// SamplerWithOverride creates a new sampler with the capability to override
// the sampling decision, if the span includes an attribute with the specified key.
// Otherwise the sampler delegates the decision to the wrapped base sampler. This
// is primarily used to enable forced tracing in Thanos components.
// Implements go.opentelemetry.io/otel/sdk/trace.Sampler interface.
func SamplerWithOverride(baseSampler tracesdk.Sampler, overrideKey attribute.Key) tracesdk.Sampler {
	return samplerWithOverride{
		baseSampler,
		overrideKey,
	}
}

func (s samplerWithOverride) ShouldSample(p tracesdk.SamplingParameters) tracesdk.SamplingResult {
	for _, attr := range p.Attributes {
		if attr.Key == s.overrideKey {
			return tracesdk.SamplingResult{
				Decision: tracesdk.RecordAndSample,
			}
		}
	}

	return s.baseSampler.ShouldSample(p)
}

func (s samplerWithOverride) Description() string {
	return fmt.Sprintf("SamplerWithOverride{%s}", string(s.overrideKey))
}
