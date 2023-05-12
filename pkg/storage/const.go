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
	"errors"
	"time"
)

var (
	ErrNoDataPoints = errors.New("no data points found") //数据不存在
	ErrNoRowsData   = errors.New("no rows given")        // row empty
	ErrUnknown      = "UNKNOWN"
)

type TimestampPrecision int

const (
	Nanoseconds TimestampPrecision = iota
	Microseconds
	Milliseconds
	Seconds
)

func (p TimestampPrecision) String() string {
	switch p {
	case Nanoseconds:
		return "ns"
	case Microseconds:
		return "us"
	case Milliseconds:
		return "ms"
	case Seconds:
		return "s"
	default:
		return ErrUnknown
	}
}

const (
	defaultPartitionDuration     = 5 * time.Minute  //数据块时间块
	defaultRetention             = 24 * time.Hour   //时间保留时间，通过checkExpiredInterval进行数据淘汰，数据最大保留时间 = defaultRetention+checkExpiredInterval
	defaultTimestampPrecision    = Seconds          //默认时间戳精度
	defaultWriteTimeout          = 30 * time.Second //数据写入超时时间
	defaultWorkersLimit          = 1                //默认处理的goroutine数
	defaultwritablePartitionsNum = 2                //默认可写入的Partition个数. 超过这时间数据丢弃
)
