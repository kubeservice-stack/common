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

package queue

import (
	"runtime"
	"sync/atomic"
)

type Item struct { //数组解构
	value interface{} //数据值
	m     bool        // mark位
}

func NewUnLockQueue(max uint64) Queue {
	q := new(UnLockQueue)
	q.maxlen = minLen(max)
	q.capM = q.maxlen - 1
	q.mp = make([]Item, q.maxlen)
	return q
}

// lock free queue
type UnLockQueue struct {
	maxlen uint64 //最大长度
	capM   uint64
	putB   uint64 //生产位
	getB   uint64 //消费位
	mp     []Item //数组
}

func (q *UnLockQueue) Maxlen() uint64 {
	return q.maxlen
}

func (q *UnLockQueue) IsEmpty() bool {
	return q.Length() == 0
}

func (q *UnLockQueue) Length() int64 {
	var putB, getB uint64
	var max uint64
	getB = q.getB
	putB = q.putB

	if putB >= getB {
		max = putB - getB
	} else {
		max = q.capM + putB - getB
	}

	return int64(max)
}

// Push queue functions
func (q *UnLockQueue) Push(val interface{}) bool {
	var putB, putBNew, getB, posCnt uint64
	var men *Item
	capM := q.capM
	for {
		getB = q.getB
		putB = q.putB

		if putB >= getB {
			posCnt = putB - getB
		} else {
			posCnt = capM + putB - getB
		}

		if posCnt >= capM {
			runtime.Gosched()
			return false
		}

		putBNew = putB + 1
		if atomic.CompareAndSwapUint64(&q.putB, putB, putBNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	men = &q.mp[putBNew&capM]

	for {
		if !men.m {
			men.value = val
			men.m = true
			return true
		} else {
			runtime.Gosched()
		}
	}
}

// get queue functions
func (q *UnLockQueue) Pop() (val interface{}, ok bool) {
	var putB, getB, getBNew, posCnt uint64
	var men *Item
	capM := q.capM
	for {
		putB = q.putB
		getB = q.getB

		if putB >= getB {
			posCnt = putB - getB
		} else {
			posCnt = capM + putB - getB
		}

		if posCnt < 1 {
			runtime.Gosched()
			return nil, false
		}

		getBNew = getB + 1
		if atomic.CompareAndSwapUint64(&q.getB, getB, getBNew) {
			break
		} else {
			runtime.Gosched()
		}
	}

	men = &q.mp[getBNew&capM]

	for {
		if men.m {
			val = men.value
			men.m = false
			return val, true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *UnLockQueue) PopMany(count int64) ([]interface{}, bool) {
	return nil, false
}

// round 到最近的2的倍数
func minLen(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}
