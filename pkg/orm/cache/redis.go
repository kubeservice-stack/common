/*
Copyright 2024 The KubeService-Stack Authors.

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
	gcache "github.com/asjdf/gorm-cache/cache"
	gconfig "github.com/asjdf/gorm-cache/config"
	gstorage "github.com/asjdf/gorm-cache/storage"
	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisCache(cfg *config.OrmCache) (gcache.Cache, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheCfg.Addr,
		Username: cfg.CacheCfg.Username,
		Password: cfg.CacheCfg.Password,
		DB:       cfg.CacheCfg.DB,
	})

	return gcache.NewGorm2Cache(&gconfig.CacheConfig{
		CacheLevel:   gconfig.CacheLevel(cfg.CacheModel.Number()),
		CacheStorage: gstorage.NewRedis(&gstorage.RedisStoreConfig{Client: redisClient}),
		// when you create/update/delete objects, invalidate cache
		InvalidateWhenUpdate:           cfg.InvalidateWhenUpdate,
		CacheTTL:                       cfg.CacheTTL,
		DisableCachePenetrationProtect: cfg.DisableCachePenetrationProtect,
		AsyncWrite:                     cfg.AsyncWrite,
	})
}
