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

func Test_GinConfig(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.GinConfig.DefaultConfig().TOML()
	assert.Equal(aa, `
[gin]
  ## APP name
  app = "application"
  ## 服务类型，支持frontend/backend, 默认backend
  server_type = "backend"
  ## 是否打开pprof
  enable_pprof = false
  ## 是否开启Health check
  enable_health = true
  ## 是否开启debug 模式
  enable_debug = false
  ## 是否开启metric接口
  enable_metrics = true
  ## 缓存开关，默认false
  enable_cache = false
  ## 是否开启签名权限验证
  enable_auth = false
  ## 是否开启指令权限验证
  enable_verify_command = false
  ## 服务启动端口
  port = 9445
  ## Trace
  trace = ""
  ## GracefulTimeout
  graceful_timeout = "3s"`)
}

func Test_GinConfigListenAddr(t *testing.T) {
	assert := assert.New(t)
	aa := GlobalCfg.ListenAddr()
	assert.Equal(aa, "0.0.0.0:9445")
}
