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

package utils

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_Md5Encode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("74b87337454200d4d33f80c4663dc5e5", Md5Encode("aaaa"), "is not equal")
}

func Test_Base64Encode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("YWFhYQ==", Base64Encode("aaaa"), "is not equal")
}

func Test_urlencode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("user%3Ddongjiang%26signature%3DeBA5HZ6lccsp1jsh%252BZ7jtDFXrR61uRHHs7RV88zc2tY%253D%26expires%3D1479390425", Urlencode("user=dongjiang&signature=eBA5HZ6lccsp1jsh%2BZ7jtDFXrR61uRHHs7RV88zc2tY%3D&expires=1479390425"))
}

func Test_urldecode(t *testing.T) {
	assert := assert.New(t)
	aa, err := Urldecode("user%3Ddongjiang%26signature%3DeBA5HZ6lccsp1jsh%252BZ7jtDFXrR61uRHHs7RV88zc2tY%253D%26expires%3D1479390425")
	assert.Equal("user=dongjiang&signature=eBA5HZ6lccsp1jsh%2BZ7jtDFXrR61uRHHs7RV88zc2tY%3D&expires=1479390425", aa)
	assert.Nil(err)
}

func Test_Uint16Encode(t *testing.T) {
	assert := assert.New(t)
	var b []byte
	aa := Uint16Encode(b, uint16(124))
	assert.NotEmpty(aa)
	assert.Equal(aa, []uint8([]byte{0x0, 0x7c}))
}

func Test_Uint16Decode(t *testing.T) {
	assert := assert.New(t)
	aa := Uint16Decode([]byte{0x0, 0x7c})
	assert.Equal(aa, uint16(124))
}
