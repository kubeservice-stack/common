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
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/kubeservice-stack/common/pkg/config"
	"gorm.io/gorm"
)

func NewSqlite3(cfg config.DBConfg) gorm.Dialector {
	return sqlite.Open(fmt.Sprintf("%s:%s/%s", cfg.Host, cfg.Port, cfg.Database))
}

func init() {
	Register(config.SQLITE3, NewSqlite3)
}
