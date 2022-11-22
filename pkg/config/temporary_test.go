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
