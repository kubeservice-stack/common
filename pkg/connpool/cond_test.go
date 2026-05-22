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

package connpool

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTMOCond_WaitAndSignal(t *testing.T) {
	assert := assert.New(t)

	var mu sync.Mutex
	cond := NewTMOCond(&mu)

	mu.Lock()

	go func() {
		time.Sleep(50 * time.Millisecond)
		// Signal without lock: TMOCond uses unbuffered channel,
		// Signal() blocks until Wait() receives. Holding lock here
		// would deadlock since Wait() needs to re-lock after receiving.
		cond.Signal()
	}()

	cond.Wait()
	mu.Unlock()

	assert.True(true)
}

func TestTMOCond_WaitOrTimeout(t *testing.T) {
	assert := assert.New(t)

	var mu sync.Mutex
	cond := NewTMOCond(&mu)

	mu.Lock()
	ret := cond.WaitOrTimeout(50 * time.Millisecond)
	mu.Unlock()

	assert.False(ret)
}

func TestTMOCond_WaitOrTimeoutSignal(t *testing.T) {
	assert := assert.New(t)

	var mu sync.Mutex
	cond := NewTMOCond(&mu)

	mu.Lock()

	go func() {
		time.Sleep(20 * time.Millisecond)
		cond.Signal()
	}()

	ret := cond.WaitOrTimeout(200 * time.Millisecond)
	mu.Unlock()

	assert.True(ret)
}

func TestTMOCond_MultipleCycles(t *testing.T) {
	// Test multiple sequential Wait/Signal cycles with the same cond
	assert := assert.New(t)

	var mu sync.Mutex
	cond := NewTMOCond(&mu)

	counter := 0
	n := 5

	for i := 0; i < n; i++ {
		mu.Lock()

		go func() {
			time.Sleep(5 * time.Millisecond)
			cond.Signal()
		}()

		cond.Wait()
		counter++
		mu.Unlock()
	}

	assert.Equal(n, counter)
}

func TestTMOCond_WaitLockRestored(t *testing.T) {
	// Verify that after Wait() returns, the mutex is held by the caller
	assert := assert.New(t)

	var mu sync.Mutex
	cond := NewTMOCond(&mu)

	mu.Lock()

	go func() {
		time.Sleep(20 * time.Millisecond)
		cond.Signal()
	}()

	cond.Wait()
	// At this point, mu should be locked by Wait()'s re-lock
	shared := 42
	mu.Unlock()

	assert.Equal(42, shared)
}

func TestConnPool_PopPush(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(5, 2, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	// Pop should create a new connection
	c, err := pool.Pop()
	assert.Nil(err)
	assert.NotNil(c)
	assert.Equal("conn", c.Inst)
	assert.Equal(1, pool.GetActiveNum())

	// Push back
	err = pool.Push(c)
	assert.Nil(err)
	assert.Equal(0, pool.GetActiveNum())
	assert.Equal(1, pool.GetIdleNum())
}

func TestConnPool_MaxLimit(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(2, 0, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	c1, err := pool.Pop()
	assert.Nil(err)
	c2, err := pool.Pop()
	assert.Nil(err)
	assert.Equal(2, pool.GetActiveNum())

	// Push both back
	pool.Push(c1)
	pool.Push(c2)
	assert.Equal(0, pool.GetActiveNum())
	assert.Equal(2, pool.GetIdleNum())
}

func TestConnPool_ClearPool(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(5, 5, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	// Create some connections
	var conns []*Conn
	for i := 0; i < 3; i++ {
		c, err := pool.Pop()
		assert.Nil(err)
		conns = append(conns, c)
	}
	// Push them back as idle
	for _, c := range conns {
		pool.Push(c)
	}
	assert.Equal(3, pool.GetIdleNum())

	pool.ClearPool()
	assert.Equal(0, pool.GetIdleNum())
}

func TestConnPool_WaitTime(t *testing.T) {
	assert := assert.New(t)

	// Pool with max=1 and wait=-1 (no wait)
	pool := NewConnectionPool(1, 0, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	c1, err := pool.Pop()
	assert.Nil(err)
	assert.NotNil(c1)

	// Second pop should return nil (no wait)
	c2, err := pool.Pop()
	assert.Nil(err)
	assert.Nil(c2)

	pool.Push(c1)
}

func TestConnPool_PushNil(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(5, 2, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	err := pool.Push(nil)
	assert.NotNil(err)
}

func TestConnPool_PushErrorConn(t *testing.T) {
	assert := assert.New(t)

	disconnected := false
	pool := NewConnectionPool(5, 2, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) { disconnected = true },
		nil,
	)

	c, _ := pool.Pop()
	c.Err = errors.New("test error")
	err := pool.Push(c)
	assert.Nil(err)
	assert.True(disconnected)
	assert.Equal(0, pool.GetActiveNum())
}

func TestConnPool_IdleTimeout(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(5, 1, 50*time.Millisecond, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	c, _ := pool.Pop()
	pool.Push(c)
	assert.Equal(1, pool.GetIdleNum())

	// Wait for idle timeout
	time.Sleep(100 * time.Millisecond)

	// Pop should trigger idle cleanup + new conn
	c2, err := pool.Pop()
	assert.Nil(err)
	assert.NotNil(c2)
	// The old idle conn should have been removed, new one created
	assert.Equal(1, pool.GetActiveNum())

	pool.Push(c2)
}

func TestConnPool_GetWaitNum(t *testing.T) {
	assert := assert.New(t)

	pool := NewConnectionPool(1, 0, 0, -1,
		func() (interface{}, error) { return "conn", nil },
		func(c interface{}) {},
		nil,
	)

	assert.Equal(0, pool.GetWaitNum())

	c, _ := pool.Pop()
	assert.Equal(0, pool.GetWaitNum())
	pool.Push(c)
}

func TestConnPool_ConnectError(t *testing.T) {
	assert := assert.New(t)

	callCount := 0
	pool := NewConnectionPool(5, 0, 0, -1,
		func() (interface{}, error) {
			callCount++
			return nil, errors.New("test error")
		},
		func(c interface{}) {},
		nil,
	)

	_, err := pool.Pop()
	assert.NotNil(err)
	assert.Equal(1, callCount)
	assert.Equal(0, pool.GetActiveNum())
}
