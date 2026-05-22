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

package cache

import (
	"time"

	"github.com/kubeservice-stack/common/pkg/cache/item"
	"github.com/kubeservice-stack/common/pkg/utils"
)

// NewARCPlugin returns a new plugin.
func NewARCPlugin(cb *Setting) Cache {
	c := &ARCPlugin{}

	options(&c.Options, cb)
	c.init()
	c.loadGroup.plugin = c

	return c
}

type ARCPlugin struct {
	Options
	items map[interface{}]*item.ArcItem

	part int
	t1   *item.ArcList
	t2   *item.ArcList
	b1   *item.ArcList
	b2   *item.ArcList
}

func (c *ARCPlugin) init() {
	c.items = make(map[interface{}]*item.ArcItem)
	c.t1 = item.NewARCList()
	c.t2 = item.NewARCList()
	c.b1 = item.NewARCList()
	c.b2 = item.NewARCList()
}

func (c *ARCPlugin) replace(key interface{}) {
	var old interface{}
	if (c.t1.Len() > 0 && c.b2.Has(key) && c.t1.Len() == c.part) || (c.t1.Len() > c.part) {
		old = c.t1.RemoveTail()
		c.b1.PushFront(old)
	} else if c.t2.Len() > 0 {
		old = c.t2.RemoveTail()
		c.b2.PushFront(old)
	} else {
		return
	}
	item, ok := c.items[old]
	if ok {
		delete(c.items, old)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(old, item.Value)
		}
	}
}

func (c *ARCPlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *ARCPlugin) set(key, value interface{}) (interface{}, error) {
	it, ok := c.items[key]
	if ok {
		it.Value = value
	} else {
		it = &item.ArcItem{
			Key:   key,
			Value: value,
		}
		c.items[key] = it
	}

	if c.expiration != nil {
		t := time.Now().Add(*c.expiration)
		it.Expiration = &t
	}

	if elt := c.b1.Lookup(key); elt != nil {
		c.part = utils.Min(c.size, c.part+utils.Max(c.b2.Len()/c.b1.Len(), 1))
		c.replace(key)
		c.b1.Remove(key, elt)
		c.t2.PushFront(key)
		return it, nil
	}

	if elt := c.b2.Lookup(key); elt != nil {
		c.part = utils.Max(0, c.part-utils.Max(c.b1.Len()/c.b2.Len(), 1))
		c.replace(key)
		c.b2.Remove(key, elt)
		c.t2.PushFront(key)
		return it, nil
	}

	if c.t1.Len()+c.b1.Len() == c.size {
		if c.t1.Len() < c.size {
			c.b1.RemoveTail()
			c.replace(key)
		} else {
			pop := c.t1.RemoveTail()
			it, ok := c.items[pop]
			if ok {
				delete(c.items, pop)
				if c.evictedFunc != nil {
					(*c.evictedFunc)(pop, it.Value)
				}
			}
		}
	} else {
		total := c.t1.Len() + c.b1.Len() + c.t2.Len() + c.b2.Len()
		if total >= c.size {
			if total == (2 * c.size) {
				c.b2.RemoveTail()
			}
			c.replace(key)
		}
	}

	c.t1.PushFront(key)

	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return it, nil
}

// Get a value from cache pool using key if it exists.
// If not exists and it has LoaderFunc, it will generate the value using you have specified LoaderFunc method returns value.
func (c *ARCPlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, true)
	}
	return v, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key, returns KeyNotFoundError.
// And send a request which refresh value for specified key if cache object has LoaderFunc.
func (c *ARCPlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *ARCPlugin) get(key interface{}) (interface{}, error) {
	rl := false
	c.mu.RLock()
	if elt := c.t1.Lookup(key); elt != nil {
		c.mu.RUnlock()
		rl = true
		c.mu.Lock()
		c.t1.Remove(key, elt)
		item := c.items[key]
		if !item.IsExpired(nil) {
			c.t2.PushFront(key)
			c.mu.Unlock()
			return item, nil
		}
		c.b2.PushFront(key)
		delete(c.items, key)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, elt.Value)
		}
		c.mu.Unlock()
	}
	if elt := c.t2.Lookup(key); elt != nil {
		c.mu.RUnlock()
		rl = true
		c.mu.Lock()
		item := c.items[key]
		if !item.IsExpired(nil) {
			c.t2.MoveToFront(elt)
			c.mu.Unlock()
			return item, nil
		}
		c.t2.Remove(key, elt)
		c.b2.PushFront(key)
		delete(c.items, key)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, elt.Value)
		}
		c.mu.Unlock()
	}

	if !rl {
		c.mu.RUnlock()
	}
	return nil, ErrCacheKeyNotFind
}

func (c *ARCPlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return it.(*item.ArcItem).Value, nil
}

func (c *ARCPlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
	if c.loaderFunc == nil {
		return nil, ErrCacheKeyNotFind
	}
	it, _, err := c.load(key, func(v interface{}, e error) (interface{}, error) {
		if e == nil {
			c.mu.Lock()
			defer c.mu.Unlock()
			return c.set(key, v)
		}
		return nil, e
	}, isWait)
	if err != nil {
		return nil, err
	}
	return it.(*item.ArcItem).Value, nil
}

// Remove removes the provided key from the cache.
func (c *ARCPlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *ARCPlugin) remove(key interface{}) bool {
	if elt := c.t1.Lookup(key); elt != nil {
		v := elt.Value
		c.t1.Remove(key, elt)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, v)
		}
		return true
	}

	if elt := c.t2.Lookup(key); elt != nil {
		v := elt.Value
		c.t2.Remove(key, elt)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, v)
		}
		return true
	}

	return false
}

func (c *ARCPlugin) keys() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	now := time.Now()
	result := make([]interface{}, 0, len(c.items))
	for _, key := range c.t1.Keys() {
		if it, ok := c.items[key]; ok && !it.IsExpired(&now) {
			result = append(result, key)
		}
	}
	for _, key := range c.t2.Keys() {
		if it, ok := c.items[key]; ok && !it.IsExpired(&now) {
			result = append(result, key)
		}
	}
	return result
}

// Keys returns a slice of the keys in the cache.
func (c *ARCPlugin) Keys() []interface{} {
	return c.keys()
}

// Returns all key-value pairs in the cache.
func (c *ARCPlugin) GetALL() map[interface{}]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	now := time.Now()
	m := make(map[interface{}]interface{})
	for _, key := range c.t1.Keys() {
		if it, ok := c.items[key]; ok && !it.IsExpired(&now) {
			m[key] = it.Value
		}
	}
	for _, key := range c.t2.Keys() {
		if it, ok := c.items[key]; ok && !it.IsExpired(&now) {
			m[key] = it.Value
		}
	}
	return m
}

// Len returns the number of non-expired items in the cache.
func (c *ARCPlugin) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	now := time.Now()
	count := 0
	for _, it := range c.items {
		if !it.IsExpired(&now) {
			count++
		}
	}
	return count
}

// Purge is used to completely clear the cache
func (c *ARCPlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.init()
}

func (c *ARCPlugin) HasKey(key interface{}) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// ARC stores items in t1/t2 lists, not just c.items
	if c.t1.Lookup(key) != nil {
		if it, ok := c.items[key]; ok {
			return !it.IsExpired(nil)
		}
	}
	if c.t2.Lookup(key) != nil {
		if it, ok := c.items[key]; ok {
			return !it.IsExpired(nil)
		}
	}
	return false
}

// init
func init() {
	Register(ARC, NewARCPlugin)
}
