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
	"time"

	"github.com/kubeservice-stack/common/pkg/logger"
)

type Option func(*Storage)

// Defaults to 5min
func WithPartitionDuration(duration time.Duration) Option {
	return func(s *Storage) {
		s.partitionDuration = duration
	}
}

// Defaults to 1d.
func WithRetention(retention time.Duration) Option {
	return func(s *Storage) {
		s.retention = retention
	}
}

// Defaults to Nanoseconds
func WithTimestampPrecision(precision TimestampPrecision) Option {
	return func(s *Storage) {
		s.timestampPrecision = precision
	}
}

// Defaults to 15s.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Storage) {
		s.writeTimeout = timeout
	}
}

// Defaults to a logger implementation that does nothing.
func WithLogger(logger *logger.Logger) Option {
	return func(s *Storage) {
		s.logger = logger
	}
}
