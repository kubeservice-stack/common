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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryLock_TryLock(t *testing.T) {
	sl, _ := NewMemoryLock()
	assert.True(t, sl.TryLock())
	assert.False(t, sl.TryLock())
}

func TestMemoryLock_Lock_Unlock(t *testing.T) {
	sl, _ := NewMemoryLock()
	sl.Lock()
	assert.False(t, sl.TryLock())

	sl.Unlock()
	assert.True(t, sl.TryLock())
	sl.Unlock()

	sl.Lock()
	go func() {
		time.Sleep(time.Millisecond)
		sl.Unlock()
	}()
	sl.Lock()
	assert.False(t, sl.TryLock())
}
