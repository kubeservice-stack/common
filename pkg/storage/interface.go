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

type StorageInterface interface {
	Reader
	// The precision of timestamps is nanoseconds by default. It can be changed using WithTimestampPrecision.
	InsertRows(rows []Row) error
	// Close gracefully shutdowns by flushing any unwritten data to the underlying disk partition.
	Close() error
}

type Reader interface {
	Select(name string, labels []Label, start, end int64) (points []*DataPoint, err error)
}

type DataPoint struct {
	// The actual value. This field must be set.
	Value float64
	// Unix timestamp.
	Timestamp int64
}

type Row struct {
	// The unique name of metric.
	Name string
	// An optional key-value properties to further detailed identification.
	Labels []Label
	// This field must be set.
	DataPoint
}
