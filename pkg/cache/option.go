package cache

import (
	"sync"
	"time"

	"github.com/kubeservice-stack/common/pkg/logger"
)

var optionsLogger = logger.GetLogger("pkg/common/cache", "option")

type Options struct {
	size        int // cache size > 0
	loaderFunc  *LoaderFunc
	evictedFunc *EvictedFunc
	addedFunc   *AddedFunc
	expiration  *time.Duration
	mu          sync.RWMutex
	loadGroup   Group
}

type LoaderFunc func(interface{}) (interface{}, error)

type EvictedFunc func(interface{}, interface{})

type AddedFunc func(interface{}, interface{})

func options(c *Options, cb *Setting) {
	c.size = cb.size
	c.loaderFunc = cb.loaderFunc
	c.expiration = cb.expiration
	c.addedFunc = cb.addedFunc
	c.evictedFunc = cb.evictedFunc
}

func (c *Options) load(key interface{}, cb func(interface{}, error) (interface{}, error), isWait bool) (interface{}, bool, error) {
	v, called, err := c.loadGroup.Do(key, func() (interface{}, error) {
		return cb((*c.loaderFunc)(key))
	}, isWait)

	if err != nil {
		optionsLogger.Error(err.Error())
		return nil, called, err
	}
	return v, called, nil
}
