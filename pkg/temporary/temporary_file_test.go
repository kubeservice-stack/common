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
	"testing"

	"github.com/kubeservice-stack/common/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var cfgpath = "./test"

func Test_NewTemporaryFile(t *testing.T) {
	assert := assert.New(t)
	err := utils.MkDirIfNotExist(cfgpath)
	assert.Nil(err)

	defer func() {
		_ = utils.RemoveDir(cfgpath)
	}()

	tf, err := newTemporaryFile(cfgpath, "test1")
	assert.Nil(err)

	defer func() {
		tf.Close()
	}()

	assert.Contains(tf.Name(), "test1")
	assert.Equal(tf.Type(), "File")
	assert.Equal(tf.Size(), int64(0))
	assert.Equal(tf.Bytes(), []byte{})

	n, err := tf.Write([]byte("dongjiang test"))
	assert.Nil(err)
	assert.Equal(tf.Size(), int64(14))
	assert.Equal(n, 14)

	n, err = tf.Read([]byte("aa"))
	assert.Nil(err)
	assert.Equal(n, 2)

	err = tf.Sync()
	assert.Nil(err)

	abs, err := tf.Seek(0, 8)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(0))

	abs, err = tf.Seek(-1, io.SeekCurrent)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Equal(abs, int64(1))
}
