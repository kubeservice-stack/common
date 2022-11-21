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
  dial_timeout = "0s"`)
}
