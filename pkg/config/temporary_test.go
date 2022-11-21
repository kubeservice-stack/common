package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Temporary(t *testing.T) {
	assert := assert.New(t)

	aa := GlobalCfg.Temporary.DefaultConfig().TOML()
	assert.Equal(aa, `
[temporary]
  ## 最大使用内存空间, 超过时则转化成文件, 默认是 5242880 byte = 5MB
  max_buffer_size = 5242880
  ## 上传文件 临时文件目录, 默认 /tmp
  file_dir = "/tmp"
  ## 上传文件临时文件名格式, 默认前缀 uploadd-*
  file_pattern = "uploadd-*"
  ## 上传文件临时文件名格式, 默认前缀 104857600 byte = 100MB
  max_upload_size = 104857600`)
}
