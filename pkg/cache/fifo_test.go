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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func evictedFuncForFIFO(key, value interface{}) {
	fmt.Printf("[FIFO] Key:%v Value:%v will evicted.\n", key, value)
}

func addFuncForFIFO(key, value interface{}) {
	fmt.Printf("[FIFO] add Key:%v Value:%v\n", key, value)
}

func optionsFIFOCache(size int, loader LoaderFunc) Cache {
	return New(size).
		FIFO().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForFIFO).
		AddedFunc(addFuncForFIFO).
		Setting()
}

func buildLoadingFIFOCache(size int, loader LoaderFunc) Cache {
	return New(size).
		FIFO().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForFIFO).
		AddedFunc(addFuncForFIFO).
		Expiration(time.Second).
		Setting()
}

func TestFIFOGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000
	gc := optionsFIFOCache(size, loader)
	// set
	for i := 0; i < numbers; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	// get
	for i := 0; i < numbers; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestFIFOGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	numbers := 1000
	gc := buildLoadingFIFOCache(size, loader)
	// set
	for i := 0; i < numbers; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	// get
	for i := 0; i < numbers; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingFIFOGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := optionsFIFOCache(size, loader)
	// get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		_, err := gc.Get(key)
		assert.NotNil(err)
	}
}

func TestLoadingFIFOGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingFIFOCache(size, loader)
	// get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		_, err := gc.Get(key)
		assert.NotNil(err)
	}
}

func TestFIFOLength(t *testing.T) {
	assert := assert.New(t)

	gc := optionsFIFOCache(1000, loader)
	gc.Set("testaa", true)
	gc.Set("testbb", false)

	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Second)

	length = gc.Len()
	assert.Equal(length, 2)
}

func TestFIFOLengthWithTimeout(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingFIFOCache(1000, loader)
	gc.Get("testaa")
	gc.Get("testbb")
	length := gc.Len()
	assert.Equal(length, 0)

	time.Sleep(time.Second)

	length = gc.Len()
	assert.Equal(length, 0)

	gc.Set("testaa", true)
	gc.Set("testbb", false)
	length = gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Second)

	length = gc.Len()
	assert.Equal(length, 2)
}

func TestFIFOEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := optionsFIFOCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.NotNil(err)
	}
}

func TestFIFOEvictItemWithTimeout(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingFIFOCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.NotNil(err)
	}
}

func TestFIFOGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		FIFO().
		LoaderFunc(
			func(key interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond)
				return "value", nil
			}).
		Setting()

	v, err := cache.GetIFPresent("key")
	assert.Equal(err, ErrCacheKeyNotFind)
	assert.Equal(v, nil)

	time.Sleep(20 * time.Millisecond) // 时间够长，case稳定

	v, err = cache.GetIFPresent("key")
	assert.Nil(err)

	assert.Equal(v, "value")
}

func TestFIFOGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		FIFO().
		Setting()

	for i := 0; i < size; i++ {
		cache.Set(i, i*i)
	}

	m := cache.GetALL()
	for i := 0; i < size; i++ {
		v, ok := m[i]
		assert.True(ok)
		assert.Equal(v, i*i)
	}
	time.Sleep(time.Millisecond)

	cache.Set(size, size*size)
	m = cache.GetALL()

	assert.Equal(len(m), 8)

	v1, ok := m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)
}

func Test_FIFONew(t *testing.T) {
	assert := assert.New(t)
	size := 8
	cache := NewFIFOPlugin(New(size).
		FIFO().
		EvictedFunc(evictedFuncForFIFO).
		Expiration(time.Millisecond))

	for i := 0; i < size; i++ {
		cache.Set(i, i*i)
	}

	ret, err := cache.Get(0)
	assert.Nil(err)
	assert.Equal(ret, 0)

	m := cache.GetALL()
	for i := 0; i < size; i++ {
		v, ok := m[i]
		assert.True(ok)
		assert.Equal(v, i*i)
	}
	r := cache.HasKey(1)
	assert.True(r)
	r = cache.Remove(1)
	assert.True(r)

	cache.Set(size, size*size)
	m = cache.GetALL()
	assert.Equal(len(m), 8)

	v1, ok := m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)

	// fmt.Println(size, m)
	v1, ok = m[0]
	assert.True(ok)
	assert.Equal(v1, 0)

	size++
	cache.Set(size, size*size)
	m = cache.GetALL()
	assert.Equal(len(m), 8)
	// fmt.Println(size, m)

	v1, ok = m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)

	cache.Purge()
}
