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

func buildARCache(size int) Cache {
	return New(size).
		ARC().
		EvictedFunc(evictedFuncForARC).
		Setting()
}

func buildLoadingARCache(size int) Cache {
	return New(size).
		ARC().
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForARC).
		Setting()
}

func buildLoadingARCacheWithExpiration(size int, ep time.Duration) Cache {
	return New(size).
		ARC().
		Expiration(ep).
		LoaderFunc(loader).
		EvictedFunc(evictedFuncForARC).
		AddedFunc(addFuncForARC).
		Setting()
}

func evictedFuncForARC(key, value interface{}) {
	fmt.Printf("[ARC] Key:%v Value:%v will evicted.\n", key, value)
}

func addFuncForARC(key, value interface{}) {
	fmt.Printf("[ARC] Add Key:%v Value:%v \n", key, value)
}

func TestARCGet(t *testing.T) {
	assert := assert.New(t)

	size := 100
	gc := buildARCache(size)

	// set
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	// get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestARCGetBig(t *testing.T) {
	assert := assert.New(t)

	size := 100
	gc := buildARCache(size)

	// set
	for i := 0; i < size+10; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
		gc.Get("Key-1")
	}

	// get
	assert.Equal(gc.Len(), size)

	for i := 0; i < 10; i++ {
		key := "Key-" + strconv.Itoa(i+size)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}

	for i := 11; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingARCGet(t *testing.T) {
	assert := assert.New(t)

	size := 100
	gc := buildLoadingARCache(size)

	// set
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
	}

	// get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestLoadingARCGetWithExpiration(t *testing.T) {
	assert := assert.New(t)

	size := 100
	gc := buildLoadingARCacheWithExpiration(size, 1*time.Nanosecond)

	// set
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		assert.Nil(err)
		gc.Set(key, value)
	}

	// get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestARCLength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingARCacheWithExpiration(2, 2*time.Second)
	gc.Set("test1", "aa")
	gc.Set("test2", "aa")
	gc.Set("test3", "aa")
	length := gc.Len()
	assert.Equal(length, 2)

	time.Sleep(time.Second * 3)
	gc.Set("test4", "aa")
	length = gc.Len()
	assert.Equal(length, 1)
}

func TestARCKeys(t *testing.T) {
	assert := assert.New(t)

	gc := buildARCache(1)
	gc.Set("test1", "aa")
	gc.Set("test2", "aa")

	ks := gc.Keys()
	assert.Equal(ks, []interface{}{"test2"})

	b := gc.Remove("test2")
	assert.True(b)

	aa := gc.HasKey("test2")
	assert.False(aa)
}

func TestARCEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingARCache(cacheSize)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.Nil(err)
	}
}

func TestARCGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		ARC().
		LoaderFunc(
			func(key interface{}) (interface{}, error) {
				return "value", nil
			}).
		Setting()

	v, err := cache.GetIFPresent("key")
	assert.Equal(err, ErrCacheKeyNotFind)

	time.Sleep(2 * time.Millisecond)

	v, err = cache.GetIFPresent("key")
	assert.Nil(err)

	assert.Equal(v, "value")
}

func TestARCGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		ARC().
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
