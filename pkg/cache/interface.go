package cache

import (
	"github.com/kubeservice-stack/common/pkg/logger"
)

type Cache interface {
	Set(interface{}, interface{})                  // set数据
	Get(interface{}) (interface{}, error)          // get数据
	GetIFPresent(interface{}) (interface{}, error) // 获取数据，如果数据不存在则通过cacheLoader获取数据，缓存并返回
	GetALL() map[interface{}]interface{}           // TODO：获得全量数据，业务慎用
	get(interface{}) (interface{}, error)          // private func: get key by serialize
	Remove(interface{}) bool                       // 删除key
	Purge()                                        // 清除 plguin
	Keys() []interface{}                           // 获得全部key
	Len() int                                      // 获得cache大小
	HasKey(interface{}) bool                       // 判断key是否存在
}

var cacheLogger = logger.GetLogger("pkg/common/cache", "interface")

type Instance func(*Setting) Cache

var adapters = make(map[MODE]Instance)

func Register(name MODE, adapter Instance) {
	if adapter == nil {
		panic("Cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func PluginInstance(cb *Setting) (adapter Cache) {
	instanceFunc, ok := adapters[cb.tp]
	if !ok {
		cacheLogger.Error("Cache: unknown adapter name %q (forgot to import?)", logger.Any("plugin", cb))
		return
	}
	adapter = instanceFunc(cb)
	return
}

func HasRegister(name MODE) bool {

	if _, ok := adapters[name]; ok {
		return true
	}
	cacheLogger.Error("Can not find adapter name: " + string(name))
	return false
}
