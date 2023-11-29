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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func loader(key interface{}) (interface{}, error) {
	return fmt.Sprintf("valueFor%s", key), nil
}

func evictedFuncForLFU(key, value interface{}) {
	fmt.Printf("[LFU] Key:%v Value:%v will evicted.\n", key, value)
}

func optionsLFUCache(size int) Cache {
	return New(size).
		LoaderFunc(loader).
		LFU().
		EvictedFunc(evictedFuncForLFU).
		Setting()
}

func optionsLoadingLFUCache(size int, loader LoaderFunc) Cache {
	return New(size).
		LFU().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForLFU).
		Expiration(time.Second).
		Setting()
}

func TestLFUGet(t *testing.T) {
	assert := assert.New(t)

	size := 100
	numbers := 100

	gc := optionsLFUCache(size)
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

func TestLFUGetWithTimeout(t *testing.T) {
	assert := assert.New(t)

	size := 100
	numbers := 100

	gc := optionsLoadingLFUCache(size, loader)
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

func TestLoadingLFUGet(t *testing.T) {
	assert := assert.New(t)

	size := 100
	numbers := 100

	gc := optionsLoadingLFUCache(size, loader)

	// get
	for i := 0; i < numbers; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLFULength(t *testing.T) {
	assert := assert.New(t)

	gc := optionsLFUCache(5)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(1 * time.Second)
	length = gc.Len()
	assert.Equal(length, 2)
}

func TestLFULengthWithTimeout(t *testing.T) {
	assert := assert.New(t)

	gc := optionsLoadingLFUCache(5, loader)
	gc.Get("test1")
	gc.Get("test2")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(1 * time.Second)
	length = gc.Len()
	assert.Equal(length, 0)
}

func TestLFUEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := optionsLFUCache(cacheSize)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.Nil(err)
	}
}

func TestLFUEvictItemWithTimeout(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := optionsLoadingLFUCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.Nil(err)
	}
}

func TestLFUGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		LFU().
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

func TestLFUGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		LFU().
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

	assert.Equal(len(m), 1)

	v1, ok := m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)
}

func Test_LFUNew(t *testing.T) {
	assert := assert.New(t)
	size := 8
	cache := NewLFUPlugin(New(size).
		LFU().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForLFU).
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
	time.Sleep(time.Millisecond * 2)

	cache.Set(size, size*size)
	m = cache.GetALL()

	assert.Equal(len(m), 1)

	v1, ok := m[size]
	assert.True(ok)
	assert.Equal(v1, size*size)
}
