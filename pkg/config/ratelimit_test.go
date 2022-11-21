package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Ratelimit(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.RateLimit.DefaultConfig().TOML()
	assert.Equal(aa, `
# 访问频率限制
[ratelimit]
  ## qps
  qps = 100
  ## 并发数
  burst = 20`)
}
