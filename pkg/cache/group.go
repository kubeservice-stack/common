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

package cache

import (
	"sync"
)

// callback type
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	plugin Cache
	mu     sync.Mutex            // protects m
	m      map[interface{}]*call // lazily initialized
}

func (g *Group) Do(key interface{}, fn func() (interface{}, error), isWait bool) (interface{}, bool, error) {
	g.mu.Lock()
	v, err := g.plugin.get(key)
	if err == nil {
		g.mu.Unlock()
		return v, false, nil
	}
	if g.m == nil {
		g.m = make(map[interface{}]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		if !isWait {
			return nil, false, ErrCacheKeyNotFind
		}
		c.wg.Wait()
		return c.val, false, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	if !isWait {
		go g.call(c, key, fn)
		return nil, false, ErrCacheKeyNotFind
	}
	v, err = g.call(c, key, fn)
	return v, true, err
}

func (g *Group) call(c *call, key interface{}, fn func() (interface{}, error)) (interface{}, error) {
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
