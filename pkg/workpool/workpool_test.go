/*
Copyright 2022 The KubeService-Stack Authors.

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

package workpool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_PoolSubmit(t *testing.T) {
	assert := assert.New(t)
	pool := NewDefaultPool("test", 2, time.Second*5)

	// Wait for dispatcher goroutine to be ready
	time.Sleep(50 * time.Millisecond)

	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.Submit(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(100)
	<-finished

	// Wait for all tasks to be processed
	assert.Eventually(func() bool {
		return c.Load() == 100
	}, time.Second*5, time.Millisecond*10)

	pool.Stop()
	pool.Stop()

	// After stop, submit 100 more tasks — all should be rejected
	go do(100)
	<-finished
	// Counter should still be 100 since pool was stopped
	assert.Equal(int32(100), c.Load())
}

func Test_PoolSubmitAndWait(t *testing.T) {
	assert := assert.New(t)
	pool := NewDefaultPool("test", 2, time.Second*5)
	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.SubmitAndWait(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(1000)
	<-finished
	assert.Equal(int32(1000), c.Load())
	pool.Stop()
}

func Test_PoolStoped(t *testing.T) {
	assert := assert.New(t)
	pool := NewDefaultPool("test", 2, time.Second*5)
	var c atomic.Int32

	finished := make(chan struct{})
	do := func(iterations int) {
		for i := 0; i < iterations; i++ {
			pool.SubmitAndWait(func() {
				c.Inc()
			})
		}
		finished <- struct{}{}
	}
	go do(1000)
	<-finished
	assert.Equal(int32(1000), c.Load())

	ret := pool.Stopped()
	assert.False(ret)
	time.Sleep(time.Second * 1)
	ret = pool.Stopped()
	assert.False(ret)
	pool.Stop()
	ret = pool.Stopped()
	assert.True(ret)
}
