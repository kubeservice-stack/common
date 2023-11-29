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

func NewSimplePlugin(cb *Setting) Cache {
	c := &SimplePlugin{}
	options(&c.Options, cb)

	c.init()
	c.loadGroup.plugin = c
	return c
}

// SimplePlugin has no clear priority for evict cache. It depends on key-value map order.
type SimplePlugin struct {
	Options
	items map[interface{}]*item.SimpleItem
}

func (c *SimplePlugin) init() {
	c.items = make(map[interface{}]*item.SimpleItem, c.size)
}

// set a new key-value pair
func (c *SimplePlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *SimplePlugin) set(key, value interface{}) (interface{}, error) {
	// Check for existing item
	it, ok := c.items[key]
	if ok {
		it.Value = value
	} else {
		// Verify size not exceeded
		if len(c.items) >= c.size {
			c.evict(1)
		}
		it = &item.SimpleItem{
			Value: value,
		}

		c.items[key] = it
	}

	if c.expiration != nil {
		t := time.Now().Add(*c.expiration)
		it.Expiration = &t
	}

	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return it, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key and has LoaderFunc,
// generate a value using `LoaderFunc` method returns value.
func (c *SimplePlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, true)
	}
	return v, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key, returns KeyNotFoundError.
// And send a request which refresh value for specified key if cache object has LoaderFunc.
func (c *SimplePlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *SimplePlugin) get(key interface{}) (interface{}, error) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if ok {
		if !item.IsExpired(nil) {
			return item, nil
		}
		c.mu.Lock()
		c.remove(key)
		c.mu.Unlock()
	}
	return nil, ErrCacheKeyNotFind
}

func (c *SimplePlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return it.(*item.SimpleItem).Value, nil
}

func (c *SimplePlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
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
	return it.(*item.SimpleItem).Value, nil
}

func (c *SimplePlugin) evict(count int) {
	now := time.Now()
	current := 0
	for key, item := range c.items {
		if current >= count {
			return
		}
		if item.Expiration == nil || now.After(*item.Expiration) {
			defer c.remove(key)
			current += 1
		}
	}
}

// Removes the provided key from the cache.
func (c *SimplePlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *SimplePlugin) remove(key interface{}) bool {
	item, ok := c.items[key]
	if ok {
		delete(c.items, key)
		if c.evictedFunc != nil {
			(*c.evictedFunc)(key, item.Value)
		}
		return true
	}
	return false
}

// Returns a slice of the keys in the cache.
func (c *SimplePlugin) keys() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]interface{}, len(c.items))
	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}
	return keys
}

// Returns a slice of the keys in the cache.
func (c *SimplePlugin) Keys() []interface{} {
	keys := []interface{}{}
	for _, k := range c.keys() {
		_, err := c.GetIFPresent(k)
		if err == nil {
			keys = append(keys, k)
		}
	}
	return keys
}

// Returns all key-value pairs in the cache.
func (c *SimplePlugin) GetALL() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	for _, k := range c.keys() {
		v, err := c.GetIFPresent(k)
		if err == nil {
			m[k] = v
		}
	}
	return m
}

// Returns the number of items in the cache.
func (c *SimplePlugin) Len() int {
	return len(c.GetALL())
}

// Completely clear the cache
func (c *SimplePlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.init()
}

func (c *SimplePlugin) HasKey(key interface{}) bool {
	return utils.InSliceIface(key, c.Keys())
}

// init
func init() {
	Register(SIMPLE, NewSimplePlugin)
}
