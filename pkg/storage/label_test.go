/*
Copyright 2023 The KubeService-Stack Authors.

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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalMetricName(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		labels []Label
		want   string
	}{
		{
			name:   "only metric",
			metric: "metric1",

			want: "metric1",
		},
		{
			name:   "missing label name",
			metric: "metric1",
			labels: []Label{
				{Value: "value1"},
			},

			want: "\x00\ametric1",
		},
		{
			name:   "missing label value",
			metric: "metric1",
			labels: []Label{
				{Key: "metric1"},
			},

			want: "\x00\ametric1",
		},
		{
			name:   "metric with a single label",
			metric: "metric1",
			labels: []Label{
				{Key: "name1", Value: "value1"},
			},
			want: "\x00\ametric1\x00\x05name1\x00\x06value1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := marshalMetricName(tt.metric, tt.labels)
			assert.Equal(t, tt.want, got)
		})
	}
}
