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
	"sync/atomic"
	"testing"
	"time"

	"github.com/kubeservice-stack/common/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestLoaderFuncLRU(t *testing.T) {
	assert := assert.New(t)

	size := 2

	testCaches := []*Setting{
		New(size).LRU(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				cacheLogger.Info("dongjiang LoaderFunc==", logger.Any("key", key))
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				cacheLogger.Info("dongjiang AddedFunc==", logger.Any("key", key), logger.Any("value", value))
			}).
			EvictedFunc(func(key, value interface{}) {
				cacheLogger.Info("dongjiang EvictedFunc==", logger.Any("key", key), logger.Any("value", value))
			}).LRU().Setting() // LRU

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()

		assert.Equal(testCounter, int64(1))
	}
}

func TestLoaderFuncLFU(t *testing.T) {
	assert := assert.New(t)

	size := 2

	testCaches := []*Setting{
		New(size).LFU(),
	}

	for _, builder := range testCaches {
		var testCounter int64
		counter := 1000
		cache := builder.
			LoaderFunc(func(key interface{}) (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				cacheLogger.Info("dongjiang LoaderFunc==", logger.Any("key", key))
				return atomic.AddInt64(&testCounter, 1), nil
			}).
			AddedFunc(func(key, value interface{}) {
				cacheLogger.Info("dongjiang AddedFunc==", logger.Any("key", key), logger.Any("value", value))
			}).
			EvictedFunc(func(key, value interface{}) {
				cacheLogger.Info("dongjiang EvictedFunc==", logger.Any("key", key), logger.Any("value", value))
			}).LFU().Setting()

		var wg sync.WaitGroup
		for i := 0; i < counter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cache.Get(0)
				assert.Nil(err)
			}()
		}
		wg.Wait()
		assert.Equal(testCounter, int64(1))
	}
}
