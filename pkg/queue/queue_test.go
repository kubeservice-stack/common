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

package queue

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPushPop(t *testing.T) {
	assert := assert.New(t)
	q := NewRingQueue(10)
	q.Push("hello")
	res, _ := q.Pop()
	assert.Equal("hello", res)
	assert.True(q.IsEmpty())
}

func TestPushPopRepeated(t *testing.T) {
	assert := assert.New(t)

	q := NewRingQueue(10)
	for i := 0; i < 100; i++ {
		q.Push("hello")
		res, _ := q.Pop()
		assert.Equal("hello", res)
		assert.True(q.IsEmpty())
	}
}

func TestPushPopMany(t *testing.T) {
	assert := assert.New(t)

	q := NewRingQueue(10)
	for i := 0; i < 10000; i++ {
		item := fmt.Sprintf("hello%v", i)
		q.Push(item)
		res, _ := q.Pop()
		assert.Equal(item, res)
	}
	assert.True(q.IsEmpty())
}

func TestPushPopMany2(t *testing.T) {
	assert := assert.New(t)

	q := NewRingQueue(10)
	for i := 0; i < 10000; i++ {
		item := fmt.Sprintf("hello%v", i)
		q.Push(item)
	}
	for i := 0; i < 10000; i++ {
		item := fmt.Sprintf("hello%v", i)
		res, _ := q.Pop()
		assert.Equal(item, res)
	}
	assert.True(q.IsEmpty())
}

func TestLfQueueConsistency(t *testing.T) {
	assert := assert.New(t)

	max := 1000000
	c := 100
	var wg sync.WaitGroup
	wg.Add(1)
	q := NewRingQueue(2)
	go func(t *testing.T) {
		i := 0
		seen := make(map[string]string)
		for {
			r, ok := q.Pop()
			if !ok {
				runtime.Gosched()

				continue
			}
			i++
			if r == nil {
				panic("consistency failure")
			}
			s := r.(string)
			_, present := seen[s]
			if present {
				t.FailNow()
				wg.Done()
				return
			}
			seen[s] = s

			if i == max {
				wg.Done()
				return
			}
		}
	}(t)

	for j := 0; j < c; j++ {
		jj := j
		cmax := max / c
		go func() {
			for i := 0; i < cmax; i++ {
				if rand.Intn(10) == 0 {
					time.Sleep(time.Duration(rand.Intn(1000)))
				}
				q.Push(fmt.Sprintf("%v %v", jj, i))
			}
		}()
	}

	wg.Wait()
	time.Sleep(500 * time.Millisecond)
	// queue should be empty
	for i := 0; i < 100; i++ {
		_, ok := q.Pop()
		assert.False(ok)
	}
}

func TestPushPopMany3(t *testing.T) {
	assert := assert.New(t)

	q := NewRingQueue(10)
	for i := 0; i < 10000; i++ {
		item := fmt.Sprintf("hello%v", i)
		q.Push(item)
	}
	for i := 0; i < 1000; i++ {
		res, ct := q.PopMany(10)
		assert.Equal(ct, true)
		assert.Equal(len(res), 10)
	}
	assert.True(q.IsEmpty())
}

func TestPushPopMany4(t *testing.T) {
	assert := assert.New(t)

	q := NewRingQueue(10)
	for i := 0; i < 10000; i++ {
		item := fmt.Sprintf("hello%v", i)
		q.Push(item)
	}
	for i := 0; i < 666; i++ {
		res, ct := q.PopMany(15)
		assert.Equal(ct, true)
		assert.Equal(len(res), 15)
	}
	res, ct := q.PopMany(15)
	assert.Equal(ct, true)
	assert.Equal(len(res), 10)
	assert.True(q.IsEmpty())
}
