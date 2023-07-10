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
	"container/list"

	"github.com/kubeservice-stack/common/pkg/cache/item"
	"github.com/kubeservice-stack/common/pkg/utils"
)

// NewFIFOPlugin returns a new plugin.
func NewFIFOPlugin(cb *Setting) Cache {
	c := &FIFOPlugin{}
	options(&c.Options, cb)

	c.init()
	c.loadGroup.plugin = c
	return c
}

type FIFOPlugin struct {
	Options
	items     map[interface{}]*list.Element
	evictList *list.List
}

func (c *FIFOPlugin) init() {
	c.evictList = list.New()
	c.items = make(map[interface{}]*list.Element, c.size+1)
}

// set a new key-value pair
func (c *FIFOPlugin) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, value)
}

func (c *FIFOPlugin) set(key, value interface{}) (interface{}, error) {
	// Check for existing item
	var it *item.FIFOItem
	if c.evictList.Len() >= c.size {
		c.evict(1)
	}
	it = &item.FIFOItem{
		Key:   key,
		Value: value,
	}
	c.items[key] = c.evictList.PushFront(it)

	if c.addedFunc != nil {
		(*c.addedFunc)(key, value)
	}

	return it, nil
}

// Get a value from cache pool using key if it exists.
// If it dose not exists key and has LoaderFunc,
// generate a value using `LoaderFunc` method returns value.
func (c *FIFOPlugin) Get(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *FIFOPlugin) GetIFPresent(key interface{}) (interface{}, error) {
	v, err := c.getValue(key)
	if err != nil {
		return c.getWithLoader(key, false)
	}
	return v, nil
}

func (c *FIFOPlugin) get(key interface{}) (interface{}, error) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil, ErrCacheKeyNotFind
	}
	return it, nil
}

func (c *FIFOPlugin) getValue(key interface{}) (interface{}, error) {
	it, err := c.get(key)
	if err != nil {
		return nil, err
	}
	return it.(*list.Element).Value.(*item.FIFOItem).Value, nil
}

func (c *FIFOPlugin) getWithLoader(key interface{}, isWait bool) (interface{}, error) {
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
	if v, ok := it.(*item.FIFOItem); ok {
		return v.Value, nil
	}
	return nil, nil
}

// evict removes the oldest item from the cache.
func (c *FIFOPlugin) evict(count int) {
	for i := 0; i < count; i++ {
		ent := c.evictList.Back()
		if ent == nil {
			return
		} else {
			c.removeElement(ent)
		}
	}
}

// Removes the provided key from the cache.
func (c *FIFOPlugin) Remove(key interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.remove(key)
}

func (c *FIFOPlugin) remove(key interface{}) bool {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

func (c *FIFOPlugin) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	entry := e.Value.(*item.FIFOItem)
	delete(c.items, entry.Key)
	if c.evictedFunc != nil {
		entry := e.Value.(*item.FIFOItem)
		(*c.evictedFunc)(entry.Key, entry.Value)
	}
}

func (c *FIFOPlugin) keys() []interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]interface{}, len(c.items))
	var i = 0
	for k := range c.items {
		keys[i] = k
		i++
	}
	return keys
}

// Returns a slice of the keys in the cache.
func (c *FIFOPlugin) Keys() []interface{} {
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
func (c *FIFOPlugin) GetALL() map[interface{}]interface{} {
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
func (c *FIFOPlugin) Len() int {
	return len(c.GetALL())
}

// Completely clear the cache
func (c *FIFOPlugin) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.init()
}

func (c *FIFOPlugin) HasKey(key interface{}) bool {
	return utils.InSliceIface(key, c.Keys())
}

// init
func init() {
	Register(FIFO, NewFIFOPlugin)
}
