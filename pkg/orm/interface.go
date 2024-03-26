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
	logging "github.com/kubeservice-stack/common/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	ErrDBNotconnected = errors.New("migrate: database not connected")
	ormLogger         = logging.GetLogger("pkg/common/orm", "orm")
)

type DBConn struct {
	config *gorm.Config
	db     *gorm.DB
}

type Instance func(cfg config.DBConfg) gorm.Dialector

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

func NewDBConn(cfg config.DBConfg) (*DBConn, error) {
	gcfg := &gorm.Config{
		DisableAutomaticPing: false,
		PrepareStmt:          false,
	}

	conn, err := gorm.Open(adapters[cfg.DBType](cfg), gcfg)
	if err != nil {
		return nil, errors.Join(ErrDBNotconnected, err)
	}

	err = conn.Use(tracing.NewPlugin())
	if err != nil {
		return nil, errors.Join(ErrDBNotconnected, err)
	}

	return &DBConn{
		config: gcfg,
		db:     conn,
	}, nil
}
