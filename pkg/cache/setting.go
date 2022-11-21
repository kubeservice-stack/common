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
