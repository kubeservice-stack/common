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

package orm

import (
	"errors"

	"github.com/kubeservice-stack/common/pkg/config"
	"github.com/kubeservice-stack/common/pkg/orm/cache"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	ErrDBNotconnected   = errors.New("migrate: database not connected")
	ErrDBTypeNotSupport = errors.New("migrate: database type not support")
)

type DBConn struct {
	config *gorm.Config
	db     *gorm.DB
}

type Instance func(cfg config.DBConfig) gorm.Dialector

var adapters = make(map[config.DBTYPE]Instance)

func Register(name config.DBTYPE, adapter Instance) {
	if adapter == nil {
		panic("orm: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("orm: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewDBConn(cfg config.DBConfig) (*DBConn, error) {
	gcfg := &gorm.Config{
		DisableAutomaticPing: false,
		PrepareStmt:          false,
	}

	adapter, ok := adapters[cfg.DBType]
	if !ok {
		return nil, ErrDBTypeNotSupport
	}

	conn, err := gorm.Open(adapter(cfg), gcfg)
	if err != nil {
		return nil, err
	}

	err = conn.Use(tracing.NewPlugin())
	if err != nil {
		return nil, err
	}

	if cfg.Cache.CacheType == config.CACHEREDIS {
		c, err := cache.NewRedisCache(&cfg.Cache)
		if err != nil {
			return nil, err
		}
		err = conn.Use(c)
		if err != nil {
			return nil, err
		}
	} else if cfg.Cache.CacheType == config.CACHEMEMORY {
		c, err := cache.NewMemoryCache(&cfg.Cache)
		if err != nil {
			return nil, err
		}
		err = conn.Use(c)
		if err != nil {
			return nil, err
		}
	}

	return &DBConn{
		config: gcfg,
		db:     conn,
	}, nil
}

func (g *DBConn) Close() error {
	db, err := g.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
