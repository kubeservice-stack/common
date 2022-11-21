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
