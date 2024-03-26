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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DBConfig(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.DBConfg.DefaultConfig().TOML()
	assert.Equal(aa, `
[database]
  ## dbtype 模式: mysql/sqlite3/postgres
  db_type = "mysql"
  ## user 用户
  user = "root"
  ## password 密码
  password = "root"
  ## Host 地址
  host = "localhost"
  ## Port 端口
  port = "3306"
  ## 数据库名
  database = "test"
  ## ca cert证书
  ca_cert = ""`)
}
