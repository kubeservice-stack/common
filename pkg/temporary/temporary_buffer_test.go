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

	"github.com/stretchr/testify/assert"
)

func Test_NewTemporaryBuffer(t *testing.T) {
	assert := assert.New(t)
	tb := newTemporaryBuffer()

	defer func() {
		tb.Close()
	}()

	assert.Contains(tb.Name(), "")
	assert.Equal(tb.Type(), "Buffer")
	assert.Equal(tb.Size(), int64(0))
	assert.Equal(tb.Bytes(), []byte(nil))

	n, err := tb.Write([]byte("dongjiang test"))
	assert.Nil(err)
	assert.Equal(tb.Size(), int64(14))
	assert.Equal(n, 14)

	n, err = tb.Read([]byte("aa"))
	assert.Nil(err)
	assert.Equal(n, 2)

	err = tb.Sync()
	assert.Nil(err)

	abs, err := tb.Seek(0, 8)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.NotNil(err)
	assert.Equal(abs, int64(0))

	abs, err = tb.Seek(-1, io.SeekCurrent)
	assert.Error(ErrBufferSeekInvalidWhence)
	assert.Nil(err)
	assert.Equal(abs, int64(1))
}
