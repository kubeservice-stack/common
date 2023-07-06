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
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMemoryQueue(t *testing.T) {
	var wg sync.WaitGroup
	var id int32

	producter := 100
	consumer := 100

	wg.Add(producter)

	q := NewUnLockQueue(1024 * 1024)

	for i := 0; i < producter; i++ {
		go func(g int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				t := fmt.Sprintf("Node.%d.%d.%d", g, j, atomic.AddInt32(&id, 1))
				q.Push(t)
			}
		}(i)
	}
	wg.Wait()

	wg.Add(consumer)
	for i := 0; i < consumer; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; {
				_, ok := q.Pop()
				if !ok {
					runtime.Gosched()
				} else {
					j++
				}
			}
		}()
	}
	wg.Wait()

	if q := q.Length(); q != 0 {
		log.Panicln("Len Error: r.len == 0", q, 0)
	} else {
		log.Println("Len:", q)
	}
}
