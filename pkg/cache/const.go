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
	LRU    MODE = "lru"    // Least Recently Used mode  最近最少使用
	LFU    MODE = "lfu"    // Least Frequently Used mode 最小频繁使用模式
	SIMPLE MODE = "simple" // Simple mode: Random 随机
	ARC    MODE = "arc"    // Adjustable Replacement Cache mode 可调换缓存模式
)
