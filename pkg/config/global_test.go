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
	"os"
	"testing"

	"github.com/caarlos0/env/v10"
	"github.com/stretchr/testify/assert"
)

func TestEnvOverrides(t *testing.T) {
	var (
		envKey   string = "GIN_APP"
		expected string = "server-override"
	)
	os.Setenv(envKey, expected)
	defer func() { os.Unsetenv(envKey) }()

	GlobalCfg = Global{
		GinConfig: GinConfig{
			App:         "server",
			EnablePprof: true,
		},
	}
	err := env.Parse(&GlobalCfg)
	assert.Equal(t, err, nil)
	assert.Equal(t, GlobalCfg.GinConfig.App, expected)
	assert.Equal(t, GlobalCfg.GinConfig.EnablePprof, true)
}

func Test_global(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.TOML()
	assert.Equal(aa, `
[logging]
  ## debug模式: stdout输出
  isterminal = false
  ## Dir是日志文件的输出目录
  dir = ""
  ## Name是日志名称
  name = ""
  ## 日志级别
  ## error, warn, info, 或者 debug
  level = ""
  ## 日志文件获取之前的最大大小（以兆字节为单位）. 默认 500MB
  maxsize = 0
  ## 要保留的最大旧日志文件数
  maxbackups = 0
  ## 根据以下情况保留旧日志文件的最大天数：时间戳编码在其文件名中； 一天定义为24小时
  maxage = 0
[metrics]
  ## flush时间周期, 默认是5秒
  flush_interval = 0
  ## 是否收集goroutine相关信息, 默认 开启为true
  enable_goruntime_metrics = false
  ## metrics_prefix, 默认前缀 application_server
  metrics_prefix = ""
  ## 自定义metric自动填充kv数据, 默认为{}
  metrics_tags = 'null'
[discovery]
  ## etcd namespace
  namespace = ""
  ## etcd 集群配置
  endpoints = []
  ## ETCD连接 timeout时间
  dial_timeout = "0s"
[gin]
  ## APP name
  app = "server-override"
  ## 服务类型，支持frontend/backend, 默认backend
  server_type = ""
  ## 是否打开pprof
  enable_pprof = true
  ## 是否开启Health check
  enable_health = false
  ## 是否开启debug 模式
  enable_debug = false
  ## 是否开启metric接口
  enable_metrics = false
  ## 缓存开关，默认false
  enable_cache = false
  ## 是否开启签名权限验证
  enable_auth = false
  ## 是否开启指令权限验证
  enable_verify_command = false
  ## 服务启动端口
  port = 0
  ## Trace
  trace = ""
  ## GracefulTimeout
  graceful_timeout = "0s"
# 访问频率限制
[ratelimit]
  ## qps
  qps = 0
  ## 并发数
  burst = 0
[temporary]
  ## 最大使用内存空间, 超过时则转化成文件, 默认是 5242880 byte = 5MB
  max_buffer_size = 0
  ## 上传文件 临时文件目录, 默认 /tmp
  file_dir = ""
  ## 上传文件临时文件名格式, 默认前缀 uploadd-*
  file_pattern = ""
  ## 上传文件临时文件名格式, 默认前缀 104857600 byte = 100MB
  max_upload_size = 0
[database]
  ## dbtype 模式: mysql/sqlite3/postgres
  db_type = ""
  ## user 用户
  user = ""
  ## password 密码
  password = ""
  ## Host 地址
  host = ""
  ## Port 端口
  port = ""
  ## 数据库名
  database = ""
  ## ca cert证书
  ca_cert = ""`)
}
