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

func Test_Discovery(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.Discovery.DefaultConfig().TOML()
	assert.Equal(aa, `
[discovery]
  ## etcd namespace
  namespace = "application"
  ## etcd 集群配置
  endpoints = ["http://127.0.0.1:2379"]
  ## ETCD连接 timeout时间
  dial_timeout = "0s"
  ## ETCD前缀key
  prefix = ""`)
}
