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
	"sync"
	"sync/atomic"
)

type ringBuffer struct {
	buffer []interface{}
	head   int64
	tail   int64
	mod    int64
}

type RingQueue struct {
	len     int64
	content *ringBuffer
	lock    sync.Mutex
}

func NewRingQueue(initialSize int64) Queue {
	return &RingQueue{
		content: &ringBuffer{
			buffer: make([]interface{}, initialSize),
			head:   0,
			tail:   0,
			mod:    initialSize,
		},
		len: 0,
	}
}

func (q *RingQueue) Push(item interface{}) bool {
	q.lock.Lock()
	c := q.content
	c.tail = (c.tail + 1) % c.mod
	if c.tail == c.head {
		var fillFactor int64 = 2
		// we need to resize

		newLen := c.mod * fillFactor
		newBuff := make([]interface{}, newLen)

		for i := int64(0); i < c.mod; i++ {
			buffIndex := (c.tail + i) % c.mod
			newBuff[i] = c.buffer[buffIndex]
		}
		newContent := &ringBuffer{
			buffer: newBuff,
			head:   0,
			tail:   c.mod,
			mod:    newLen,
		}
		q.content = newContent
	}
	atomic.AddInt64(&q.len, 1)
	q.content.buffer[q.content.tail] = item
	q.lock.Unlock()
	return true
}

func (q *RingQueue) Length() int64 {
	return atomic.LoadInt64(&q.len)
}

func (q *RingQueue) IsEmpty() bool {
	return q.Length() == 0
}

func (q *RingQueue) Pop() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	// as we are a single consumer, no other thread can have poped the items there are guaranteed to be items now

	q.lock.Lock()
	c := q.content
	c.head = (c.head + 1) % c.mod
	res := c.buffer[c.head]
	c.buffer[c.head] = nil
	atomic.AddInt64(&q.len, -1)
	q.lock.Unlock()
	return res, true
}

func (q *RingQueue) PopMany(count int64) ([]interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}

	q.lock.Lock()
	c := q.content

	if count >= q.len {
		count = q.len
	}
	atomic.AddInt64(&q.len, -count)

	buffer := make([]interface{}, count)
	for i := int64(0); i < count; i++ {
		pos := (c.head + 1 + i) % c.mod
		buffer[i] = c.buffer[pos]
		c.buffer[pos] = nil
	}
	c.head = (c.head + count) % c.mod

	q.lock.Unlock()
	return buffer, true
}
