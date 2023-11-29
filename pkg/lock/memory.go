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

package lock

import (
	"runtime"
	"sync/atomic"
)

// MemoryLock implements sync/Locker, default 0 indicates an unlocked memory.
type MemoryLock struct {
	_flag uint32
}

// NewMemoryLock create new memory lock instance
func NewMemoryLock() (Locker, error) {
	return &MemoryLock{
		_flag: 0,
	}, nil
}

// Lock locks memory. If the lock is locked before, the caller will be blocked until unlocked.
func (sl *MemoryLock) Lock() error {
	for !sl.TryLock() {
		runtime.Gosched() // allow other goroutines to do stuff.
	}
	return nil
}

// Unlock unlocks memory, this operation is reentrant。
func (sl *MemoryLock) Unlock() error {
	atomic.StoreUint32(&sl._flag, 0)
	return nil
}

// TryLock will try to lock memory and return whether it succeed or not without blocking.
func (sl *MemoryLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl._flag, 0, 1)
}
