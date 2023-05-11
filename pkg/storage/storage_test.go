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

func Test_storage_Select(t *testing.T) {
	tests := []struct {
		name    string
		storage Storage
		metric  string
		labels  []Label
		start   int64
		end     int64
		want    []*DataPoint
		wantErr bool
	}{
		{
			name:   "select from single partition",
			metric: "metric1",
			start:  1,
			end:    4,
			storage: func() Storage {
				part1 := NewMemoryPartition(1*time.Hour, Seconds)
				_, err := part1.insertRows([]Row{
					{DataPoint: DataPoint{Timestamp: 1}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 2}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 3}, Name: "metric1"},
				})
				if err != nil {
					panic(err)
				}
				list := newPartitionList()
				list.insert(part1)
				return Storage{
					partitionList:  list,
					workersLimitCh: make(chan struct{}, defaultWorkersLimit),
				}
			}(),
			want: []*DataPoint{
				{Timestamp: 1},
				{Timestamp: 2},
				{Timestamp: 3},
			},
		},
		{
			name:   "select from three partitions",
			metric: "metric1",
			start:  1,
			end:    10,
			storage: func() Storage {
				part1 := NewMemoryPartition(1*time.Hour, Seconds)
				_, err := part1.insertRows([]Row{
					{DataPoint: DataPoint{Timestamp: 1}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 2}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 3}, Name: "metric1"},
				})
				if err != nil {
					panic(err)
				}
				part2 := NewMemoryPartition(1*time.Hour, Seconds)
				_, err = part2.insertRows([]Row{
					{DataPoint: DataPoint{Timestamp: 4}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 5}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 6}, Name: "metric1"},
				})
				if err != nil {
					panic(err)
				}
				part3 := NewMemoryPartition(1*time.Hour, Seconds)
				_, err = part3.insertRows([]Row{
					{DataPoint: DataPoint{Timestamp: 7}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 8}, Name: "metric1"},
					{DataPoint: DataPoint{Timestamp: 9}, Name: "metric1"},
				})
				if err != nil {
					panic(err)
				}
				list := newPartitionList()
				list.insert(part1)
				list.insert(part2)
				list.insert(part3)

				return Storage{
					partitionList:  list,
					workersLimitCh: make(chan struct{}, defaultWorkersLimit),
				}
			}(),
			want: []*DataPoint{
				{Timestamp: 1},
				{Timestamp: 2},
				{Timestamp: 3},
				{Timestamp: 4},
				{Timestamp: 5},
				{Timestamp: 6},
				{Timestamp: 7},
				{Timestamp: 8},
				{Timestamp: 9},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.storage.Select(tt.metric, tt.labels, tt.start, tt.end)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
