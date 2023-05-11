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
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_memoryPartition_InsertRows(t *testing.T) {
	tests := []struct {
		name               string
		memoryPartition    *memoryPartition
		rows               []Row
		wantErr            bool
		wantDataPoints     []*DataPoint
		wantOutOfOrderRows []Row
	}{
		{
			name:            "insert in-order rows",
			memoryPartition: NewMemoryPartition(0, Seconds).(*memoryPartition),
			rows: []Row{
				{Name: "metric1", DataPoint: DataPoint{Timestamp: 1, Value: 0.1}},
				{Name: "metric1", DataPoint: DataPoint{Timestamp: 2, Value: 0.1}},
				{Name: "metric1", DataPoint: DataPoint{Timestamp: 3, Value: 0.1}},
			},
			wantDataPoints: []*DataPoint{
				{Timestamp: 1, Value: 0.1},
				{Timestamp: 2, Value: 0.1},
				{Timestamp: 3, Value: 0.1},
			},
			wantOutOfOrderRows: []Row{},
		},
		{
			name: "insert out-of-order rows",
			memoryPartition: func() *memoryPartition {
				m := NewMemoryPartition(0, Microseconds).(*memoryPartition)
				m.insertRows([]Row{
					{Name: "metric1", DataPoint: DataPoint{Timestamp: 2, Value: 0.1}},
				})
				return m
			}(),
			rows: []Row{
				{Name: "metric1", DataPoint: DataPoint{Timestamp: 1, Value: 0.1}},
			},
			wantDataPoints: []*DataPoint{
				{Timestamp: 2, Value: 0.1},
			},
			wantOutOfOrderRows: []Row{
				{Name: "metric1", DataPoint: DataPoint{Timestamp: 1, Value: 0.1}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutOfOrder, err := tt.memoryPartition.insertRows(tt.rows)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantOutOfOrderRows, gotOutOfOrder)

			got, _ := tt.memoryPartition.selectDataPoints("metric1", nil, 0, 4)
			assert.Equal(t, tt.wantDataPoints, got)
		})
	}
}

func Test_memoryPartition_SelectDataPoints(t *testing.T) {
	tests := []struct {
		name            string
		metric          string
		labels          []Label
		start           int64
		end             int64
		memoryPartition *memoryPartition
		want            []*DataPoint
	}{
		{
			name:            "given non-exist metric name",
			metric:          "unknown",
			start:           1,
			end:             2,
			memoryPartition: NewMemoryPartition(0, Milliseconds).(*memoryPartition),
			want:            []*DataPoint{},
		},
		{
			name:   "select some points",
			metric: "metric1",
			start:  2,
			end:    4,
			memoryPartition: func() *memoryPartition {
				m := NewMemoryPartition(0, Microseconds).(*memoryPartition)
				m.insertRows([]Row{
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 1, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 2, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 3, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 4, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 5, Value: 0.1},
					},
				})
				return m
			}(),
			want: []*DataPoint{
				{Timestamp: 2, Value: 0.1},
				{Timestamp: 3, Value: 0.1},
			},
		},
		{
			name:   "select all points",
			metric: "metric1",
			start:  1,
			end:    4,
			memoryPartition: func() *memoryPartition {
				m := NewMemoryPartition(0, Nanoseconds).(*memoryPartition)
				m.insertRows([]Row{
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 1, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 2, Value: 0.1},
					},
					{
						Name:      "metric1",
						DataPoint: DataPoint{Timestamp: 3, Value: 0.1},
					},
				})
				return m
			}(),
			want: []*DataPoint{
				{Timestamp: 1, Value: 0.1},
				{Timestamp: 2, Value: 0.1},
				{Timestamp: 3, Value: 0.1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.memoryPartition.selectDataPoints(tt.metric, tt.labels, tt.start, tt.end)
			assert.Equal(t, tt.want, got)
		})
	}
}

type fakeEncoder struct {
	encodePointFunc func(*DataPoint) error
	flushFunc       func() error
}

func (f *fakeEncoder) encodePoint(p *DataPoint) error {
	if f.encodePointFunc == nil {
		return nil
	}
	return f.encodePointFunc(p)
}

func (f *fakeEncoder) flush() error {
	if f.flushFunc == nil {
		return nil
	}
	return f.flushFunc()
}

func Test_toUnix(t *testing.T) {
	tests := []struct {
		name      string
		t         time.Time
		precision TimestampPrecision
		want      int64
	}{
		{
			name:      "to nanosecond",
			t:         time.Unix(1600000000, 0),
			precision: Nanoseconds,
			want:      1600000000000000000,
		},
		{
			name:      "to microsecond",
			t:         time.Unix(1600000000, 0),
			precision: Microseconds,
			want:      1600000000000000,
		},
		{
			name:      "to millisecond",
			t:         time.Unix(1600000000, 0),
			precision: Milliseconds,
			want:      1600000000000,
		},
		{
			name:      "to second",
			t:         time.Unix(1600000000, 0),
			precision: Seconds,
			want:      1600000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toUnix(tt.t, tt.precision)
			assert.Equal(t, tt.want, got)
		})
	}
}
