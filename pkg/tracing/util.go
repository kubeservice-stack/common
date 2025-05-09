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
	"testing"

	"github.com/efficientgo/core/testutil"

	opentracing "github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func CountSpans_ClientEnablesTracing(t *testing.T, exp *tracetest.InMemoryExporter, clientRoot, srvRoot, srvChild opentracing.Span) {
	testutil.Equals(t, 0, len(exp.GetSpans()))

	srvChild.Finish()
	testutil.Equals(t, 1, len(exp.GetSpans()))
	testutil.Equals(t, 1, CountSampledSpans(exp.GetSpans()))

	srvRoot.Finish()
	testutil.Equals(t, 2, len(exp.GetSpans()))
	testutil.Equals(t, 2, CountSampledSpans(exp.GetSpans()))

	clientRoot.Finish()
	testutil.Equals(t, 3, len(exp.GetSpans()))
	testutil.Equals(t, 3, CountSampledSpans(exp.GetSpans()))
}

func ContextTracing_ClientDisablesTracing(t *testing.T, exp *tracetest.InMemoryExporter, clientRoot, srvRoot, srvChild opentracing.Span) {
	testutil.Equals(t, 0, len(exp.GetSpans()))

	// Since we are not recording neither sampling, no spans should show up.
	srvChild.Finish()
	testutil.Equals(t, 0, len(exp.GetSpans()))

	srvRoot.Finish()
	testutil.Equals(t, 0, len(exp.GetSpans()))

	clientRoot.Finish()
	testutil.Equals(t, 0, len(exp.GetSpans()))
}

func ContextTracing_ForceTracing(t *testing.T, exp *tracetest.InMemoryExporter, clientRoot, srvRoot, srvChild opentracing.Span) {
	testutil.Equals(t, 0, len(exp.GetSpans()))

	srvChild.Finish()
	testutil.Equals(t, 1, len(exp.GetSpans()))
	testutil.Equals(t, 1, CountSampledSpans(exp.GetSpans()))

	srvRoot.Finish()
	testutil.Equals(t, 2, len(exp.GetSpans()))
	testutil.Equals(t, 2, CountSampledSpans(exp.GetSpans()))

	clientRoot.Finish()
	testutil.Equals(t, 3, len(exp.GetSpans()))
	testutil.Equals(t, 3, CountSampledSpans(exp.GetSpans()))
}

// Utility function for use with tests in pkg/tracing.
func CountSampledSpans(ss tracetest.SpanStubs) int {
	var count int
	for _, s := range ss {
		if s.SpanContext.IsSampled() {
			count++
		}
	}

	return count
}
