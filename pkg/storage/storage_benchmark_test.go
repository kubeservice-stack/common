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

package storage_test

import (
	"testing"

	"github.com/kubeservice-stack/common/pkg/storage"
	"github.com/stretchr/testify/require"
)

func BenchmarkStorage_InsertRows(b *testing.B) {
	stg, err := storage.NewStorage()
	require.NoError(b, err)
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		stg.InsertRows([]storage.Row{
			storage.Row{Name: "metric1", DataPoint: storage.DataPoint{Timestamp: int64(i), Value: 0.1}},
		})
	}
}

// Select data points among a thousand data in memory
func BenchmarkStorage_SelectAmongThousandPoints(b *testing.B) {
	stg, err := storage.NewStorage()
	require.NoError(b, err)
	for i := 1; i < 1000; i++ {
		stg.InsertRows([]storage.Row{
			{Name: "metric1", DataPoint: storage.DataPoint{Timestamp: int64(i), Value: 0.1}},
		})
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = stg.Select("metric1", nil, 10, 100)
	}
}

// Select data points among a million data in memory
func BenchmarkStorage_SelectAmongMillionPoints(b *testing.B) {
	stg, err := storage.NewStorage()
	require.NoError(b, err)
	for i := 1; i < 1000000; i++ {
		stg.InsertRows([]storage.Row{
			{Name: "metric1", DataPoint: storage.DataPoint{Timestamp: int64(i), Value: 0.1}},
		})
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = stg.Select("metric1", nil, 10, 100)
	}
}
