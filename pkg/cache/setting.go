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
	"time"

	"github.com/kubeservice-stack/common/pkg/logger"
)

var settingLogger = logger.GetLogger("pkg/common/cache", "setting")

func New(size int) *Setting {
	if size <= 0 {
		settingLogger.Error(ErrCacheSizeZero.Error())
		panic(ErrCacheSizeZero.Error())
	}

	return &Setting{
		tp:   LFU, // 默认LFU
		size: size,
	}
}

type Setting struct {
	tp          MODE // mode : lru \ lfu
	size        int  // cache size > 0
	loaderFunc  *LoaderFunc
	evictedFunc *EvictedFunc
	addedFunc   *AddedFunc
	expiration  *time.Duration
}

func (cb *Setting) LoaderFunc(loaderFunc LoaderFunc) *Setting {
	cb.loaderFunc = &loaderFunc
	return cb
}

func (cb *Setting) EvictType(tp MODE) *Setting {
	cb.tp = tp
	return cb
}

func (cb *Setting) LRU() *Setting {
	return cb.EvictType(LRU)
}

func (cb *Setting) LFU() *Setting {
	return cb.EvictType(LFU)
}

func (cb *Setting) Simple() *Setting {
	return cb.EvictType(SIMPLE)
}

func (cb *Setting) ARC() *Setting {
	return cb.EvictType(ARC)
}

func (cb *Setting) EvictedFunc(evictedFunc EvictedFunc) *Setting {
	cb.evictedFunc = &evictedFunc
	return cb
}

func (cb *Setting) AddedFunc(addedFunc AddedFunc) *Setting {
	cb.addedFunc = &addedFunc
	return cb
}

func (cb *Setting) Expiration(expiration time.Duration) *Setting {
	cb.expiration = &expiration
	return cb
}

func (cb *Setting) Setting() Cache {
	if HasRegister(cb.tp) {
		return PluginInstance(cb)
	} else {
		settingLogger.Error(ErrCacheUnknownAdapter.Error() + string(cb.tp))
		panic(ErrCacheUnknownAdapter.Error() + string(cb.tp))
	}
}
