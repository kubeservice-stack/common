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
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func buildSimpleCache(size int) Cache {
	return New(size).
		Simple().
		EvictedFunc(evictedFuncForSimple).
		Setting()
}

func buildLoadingSimpleCache(size int, loader LoaderFunc) Cache {
	return New(size).
		LoaderFunc(loader).
		Simple().
		EvictedFunc(evictedFuncForSimple).
		AddedFunc(addFuncForSimple).
		Setting()
}

func evictedFuncForSimple(key, value interface{}) {
	fmt.Printf("[Simple] Key:%v Value:%v will evicted.\n", key, value)
}

func addFuncForSimple(key, value interface{}) {
	fmt.Printf("[Simple] Add Key:%v Value:%v \n", key, value)
}

func TestSimpleGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildSimpleCache(size)

	//set
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestSimpleGetBig(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildSimpleCache(size)

	//set
	for i := 0; i < size+10; i++ {
		key := "Key-" + strconv.Itoa(i)
		value, err := loader(key)
		if err != nil {
			t.Error(err)
			return
		}
		gc.Set(key, value)
	}

	//get
	assert.Equal(gc.Len(), size)
}

func TestLoadingSimpleGet(t *testing.T) {
	assert := assert.New(t)

	size := 1000
	gc := buildLoadingSimpleCache(size, loader)

	//get
	for i := 0; i < size; i++ {
		key := "Key-" + strconv.Itoa(i)
		v, err := gc.Get(key)
		assert.Nil(err)
		expectedV, _ := loader(key)
		assert.Equal(v, expectedV)
	}
}

func TestSimpleLength(t *testing.T) {
	assert := assert.New(t)

	gc := buildLoadingSimpleCache(1000, loader)
	gc.Get("test1")
	gc.Get("test2")

	length := gc.Len()
	expectedLength := 2
	assert.Equal(length, expectedLength)
	log.Println("dongjiang223", gc.GetALL())
}

func TestSimpleLength2(t *testing.T) {
	assert := assert.New(t)

	gc := buildSimpleCache(1000)
	gc.Get("test1")
	gc.Get("test2")

	length := gc.Len()
	expectedLength := 0
	assert.Equal(length, expectedLength)
	log.Println("dongjiang123", gc.GetALL())
}

func TestSimpleKeys(t *testing.T) {
	assert := assert.New(t)

	gc := buildSimpleCache(1)
	gc.Set("test1", "aa")
	gc.Set("test2", "aa")

	ks := gc.Keys()
	assert.Equal(ks, []interface{}{"test2"})

	b := gc.Remove("test2")
	assert.True(b)

	aa := gc.HasKey("test2")
	assert.False(aa)
}

func TestSimpleEvictItem(t *testing.T) {
	assert := assert.New(t)

	cacheSize := 10
	numbers := 11
	gc := buildLoadingSimpleCache(cacheSize, loader)

	for i := 0; i < numbers; i++ {
		_, err := gc.Get("Key-" + strconv.Itoa(i))
		assert.Nil(err)
	}
}

func TestSimpleGetIFPresent(t *testing.T) {
	assert := assert.New(t)

	cache := New(8).
		EvictType(SIMPLE).
		LoaderFunc(
			func(key interface{}) (interface{}, error) {
				time.Sleep(time.Millisecond)
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

func TestSimpleGetALL(t *testing.T) {
	assert := assert.New(t)

	size := 8
	cache := New(size).
		Expiration(time.Millisecond).
		EvictType(SIMPLE).
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
