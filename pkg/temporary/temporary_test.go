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

package temporary

import (
	"io"
	"strings"
	"testing"

	"github.com/kubeservice-stack/common/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_NewTemporary(t *testing.T) {
	assert := assert.New(t)
	err := utils.MkDirIfNotExist(cfgpath)
	assert.Nil(err)

	defer func() {
		_ = utils.RemoveDir(cfgpath)
	}()

	sb := new(strings.Reader)
	aa, err := NewTemporary(sb, 256, cfgpath, "*")
	assert.Nil(err)

	defer aa.Close()

	assert.Equal(aa.Name(), "")
	assert.Equal(aa.Bytes(), []byte(nil))
	assert.Equal(aa.Size(), int64(0))
	assert.Equal(aa.Type(), "Buffer")

	n, err := aa.Read([]byte("aa"))
	assert.NotNil(err)
	assert.Equal(n, 0)

	abs, err := aa.Seek(0, 8)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(0))

	abs, err = aa.Seek(-1, io.SeekCurrent)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(-1))
}

func Test_NewAsyncTemporary(t *testing.T) {
	assert := assert.New(t)
	err := utils.MkDirIfNotExist(cfgpath)
	assert.Nil(err)

	defer func() {
		_ = utils.RemoveDir(cfgpath)
	}()

	sb := new(strings.Reader)
	aa := NewAsyncTemporary(sb, 256, cfgpath, "*")

	defer aa.Close()

	assert.Equal(aa.Name(), "")
	assert.Equal(aa.Bytes(), []byte(nil))
	assert.Equal(aa.Size(), int64(0))
	assert.Equal(aa.Type(), "Buffer")

	n, err := aa.Read([]byte("aa"))
	assert.NotNil(err)
	assert.Equal(n, 0)

	abs, err := aa.Seek(0, 8)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(0))

	abs, err = aa.Seek(-1, io.SeekCurrent)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(-1))
}
