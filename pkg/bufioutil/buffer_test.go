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

package bufioutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	buffer := NewBuffer([]byte{1, 2, 3})
	b, err := buffer.GetByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), b)

	buffer.SetIdx(100)
	_, err = buffer.GetByte()
	assert.Equal(t, errOutOfRange, err)

	// reset
	buffer.SetBuf([]byte{1, 2, 3})
	buffer.SetIdx(1)
	b, err = buffer.GetByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(2), b)
}
