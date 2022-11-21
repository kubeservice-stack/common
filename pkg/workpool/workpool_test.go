package workpool

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func Test_PoolSubmit(t *testing.T) {
	assert := assert.New(t)
	grNum := runtime.NumGoroutine()
	pool := NewDefaultPool("test", 2, time.Second*5)
	// 1个dispatcher goroutine + N个 work goroutine
	assert.Equal(grNum+1, runtime.NumGoroutine())

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
	assert.True(grNum+2+1 <= runtime.NumGoroutine())
	pool.Stop()
	pool.Stop()
	// reject all task
	go do(100)
	<-finished
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
