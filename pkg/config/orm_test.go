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

	aa := GlobalCfg.DBConfig.DefaultConfig().TOML()
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
  ## args 参数
  args = ""
[dbcache]
  ## orm_cache_type: redis, memory. 默认值为memory
  orm_cache_type: memory
  ## orm_cache_model 模式: CacheDisable, CacheOnlyPrimary, CacheOnlySearch 和 CacheAll. 默认值为CacheOnlySearch
  orm_cache_model = CacheOnlySearch
  ## orm_invalidate_when_update: 如果用户在数据库中更新/删除/创建某些内容，我们将使所有缓存数据无效以确保一致性. 如果true将会清理cache. 默认值为true
  orm_invalidate_when_update = true
  ## orm_async_write: 如果 true，那么我们将以异步模式写入缓存. 默认值为true
  orm_async_write = true
  ## orm_cache_ttl: CacheTTL 缓存 ttl，单位为 ms，其中 0 代表永远. 默认值为5000
  orm_cache_ttl = 5000
  ## orm_disable_cache_penetration_protect: 如果为 true，那么我们不会缓存 nil 结果, 默认值为false
  orm_disable_cache_penetration_protect = false
[dbcachecfg]
  ## cache addr
  orm_cache_addr = 
  ## cache username
  orm_cache_username = 
  ## cache password
  orm_cache_password = 
  ## cache db id
  orm_cache_dbid = 0`)
}
