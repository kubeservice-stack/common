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

package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	name      string
	capacity  uint64
	allowance uint64
	max       uint64
	unit      uint64
	lastCheck uint64
}

// New 创建 TokenBucket 实例
func New(name string, capacity uint64, tw time.Duration) *TokenBucket {
	nano := uint64(tw)
	if nano < 1 {
		nano = uint64(time.Second)
	}
	if capacity < 1 {
		capacity = 1
	}

	return &TokenBucket{
		name:      name,
		capacity:  uint64(capacity),
		allowance: uint64(capacity) * nano,
		max:       uint64(capacity) * nano,
		unit:      nano,
		lastCheck: unixNano(),
	}
}

// Limit 判断是否超过限制
func (rl *TokenBucket) Limit() bool {
	now := unixNano()
	// 计算上一次调用到现在过了多少纳秒
	passed := now - atomic.SwapUint64(&rl.lastCheck, now)

	capacity := atomic.LoadUint64(&rl.capacity)
	current := atomic.AddUint64(&rl.allowance, passed*capacity)

	if max := atomic.LoadUint64(&rl.max); current > max {
		atomic.AddUint64(&rl.allowance, max-current)
		current = max
	}

	if current < rl.unit {
		return true
	}

	// 没有超过限额
	atomic.AddUint64(&rl.allowance, -rl.unit)
	return false
}

// UpdateRate 更新速率值
func (rl *TokenBucket) UpdateCapacity(capacity uint64) {
	atomic.StoreUint64(&rl.capacity, capacity)
	atomic.StoreUint64(&rl.max, capacity*rl.unit)
}

// Undo 重置上一次调用Limit()，返回没有使用过的限额
func (rl *TokenBucket) Undo() {
	current := atomic.AddUint64(&rl.allowance, rl.unit)
	max := atomic.LoadUint64(&rl.max)
	if current > max {
		atomic.AddUint64(&rl.allowance, max-current)
	}
}

// unixNano 当前时间（纳秒）
func unixNano() uint64 {
	return uint64(time.Now().UnixNano())
}
