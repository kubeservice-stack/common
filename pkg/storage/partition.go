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

type partition interface {
	// Write operations

	insertRows(rows []Row) (outdatedRows []Row, err error)
	clean() error

	// Read operations

	selectDataPoints(metric string, labels []Label, start, end int64) ([]*DataPoint, error)
	// minTimestamp returns the minimum Unix timestamp in milliseconds.
	minTimestamp() int64
	// maxTimestamp returns the maximum Unix timestamp in milliseconds.
	maxTimestamp() int64
	// size returns the number of data points the partition holds.
	size() int
	// active means not only writable but having the qualities to be the head partition.
	active() bool
	// expired means it should get removed.
	expired() bool
}
