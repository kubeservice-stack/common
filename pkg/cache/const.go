package cache

import (
	"fmt"
)

var (
	ErrCacheSizeZero           = fmt.Errorf("Cache: Size <= 0")
	ErrCacheRegisterAdapterNil = fmt.Errorf("Cache: Register adapter is nil")
	ErrCacheCanNotFindAdapter  = fmt.Errorf("Cache: Can not find adapter: ")
	ErrCacheUnknownAdapter     = fmt.Errorf("Cache: unknown adapter: ")
	ErrCacheKeyNotFind         = fmt.Errorf("Cache: key not find")
)

type MODE string

const (
	LRU MODE = "lru" // Least Recently Used mode  最近最少使用
	LFU MODE = "lfu" // Least Frequently Used mode 最小频繁使用模式
)
