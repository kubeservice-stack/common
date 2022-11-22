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

func Test_Metrics(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.Metrics.DefaultConfig().TOML()
	assert.Equal(aa, `
[metrics]
  ## flush时间周期, 默认是5秒
  flush_interval = 5
  ## 是否收集goroutine相关信息, 默认 开启为true
  enable_goruntime_metrics = true
  ## metrics_prefix, 默认前缀 application_server
  metrics_prefix = "application_server"
  ## 自定义metric自动填充kv数据, 默认为{}
  metrics_tags = '{}'`)
}
