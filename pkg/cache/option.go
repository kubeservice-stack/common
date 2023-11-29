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
