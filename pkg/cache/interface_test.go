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

	var testCaches = []*Setting{
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

	var testCaches = []*Setting{
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
