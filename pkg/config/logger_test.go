package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_logger(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.Logging.DefaultConfig().TOML()
	assert.Equal(aa, `
[logging]
  ## debug模式: stdout输出
  isterminal = false
  ## Dir是日志文件的输出目录
  dir = "/tmp/media/log"
  ## Name是日志名称
  name = "media.log"
  ## 日志级别
  ## error, warn, info, 或者 debug
  level = "info"
  ## 日志文件获取之前的最大大小（以兆字节为单位）. 默认 500MB
  maxsize = 500
  ## 要保留的最大旧日志文件数
  maxbackups = 10
  ## 根据以下情况保留旧日志文件的最大天数：时间戳编码在其文件名中； 一天定义为24小时
  maxage = 30`)
}
