/*
Copyright 2023 The KubeService-Stack Authors.

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

package mcpack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt8(t *testing.T) {
	assert := assert.New(t)
	aa := Int8([]byte("a"))
	assert.Equal(aa, int8(97))
}

func TestPutInt8(t *testing.T) {
	assert := assert.New(t)
	a := []byte(" ")
	PutInt8(a, int8(97))
	assert.Equal(a, []byte("a"))
}

func TestInt16(t *testing.T) {
	assert := assert.New(t)
	aa := Int16([]byte("aa"))
	assert.Equal(aa, int16(24929))
}

func TestPutInt16(t *testing.T) {
	assert := assert.New(t)
	a := []byte("tt")
	PutInt16(a, int16(24929))
	assert.Equal(a, []byte("aa"))
}

func TestInt32(t *testing.T) {
	assert := assert.New(t)
	aa := Int32([]byte("aaaa"))
	assert.Equal(aa, int32(1633771873))
}

func TestPutInt32(t *testing.T) {
	assert := assert.New(t)
	a := []byte("tttt")
	PutInt32(a, int32(1633771873))
	assert.Equal(a, []byte("aaaa"))
}

func TestInt64(t *testing.T) {
	assert := assert.New(t)
	aa := Int64([]byte("dongjiang"))
	assert.Equal(aa, int64(7953754322635747172))
}

func TestPutInt64(t *testing.T) {
	assert := assert.New(t)
	a := []byte("tttttttt")
	PutInt64(a, int64(7953754322635747172))
	assert.Equal(a, []byte("dongjian"))
}

func TestUInt8(t *testing.T) {
	assert := assert.New(t)
	aa := Uint8([]byte("a"))
	assert.Equal(aa, uint8(97))
}

func TestUPutInt8(t *testing.T) {
	assert := assert.New(t)
	a := []byte(" ")
	PutUint8(a, uint8(97))
	assert.Equal(a, []byte("a"))
}

func TestUInt16(t *testing.T) {
	assert := assert.New(t)
	aa := Uint16([]byte("aa"))
	assert.Equal(aa, uint16(0x6161))
}

func TestUPutInt16(t *testing.T) {
	assert := assert.New(t)
	a := []byte("  ")
	PutUint16(a, uint16(0x6161))
	assert.Equal(a, []byte("aa"))
}

func TestUInt32(t *testing.T) {
	assert := assert.New(t)
	aa := Uint32([]byte("aaaa"))
	assert.Equal(aa, uint32(0x61616161))
}

func TestUPutInt32(t *testing.T) {
	assert := assert.New(t)
	a := []byte("  rr")
	PutUint32(a, uint32(0x61616161))
	assert.Equal(a, []byte("aaaa"))
}

func TestUInt64(t *testing.T) {
	assert := assert.New(t)
	aa := Uint64([]byte("aaaaaaaa"))
	assert.Equal(aa, uint64(0x6161616161616161))
}

func TestUPutInt64(t *testing.T) {
	assert := assert.New(t)
	a := []byte("  rrwwww")
	PutUint64(a, uint64(0x6161616161616161))
	assert.Equal(a, []byte("aaaaaaaa"))
}

func TestFloat32(t *testing.T) {
	assert := assert.New(t)
	aa := Float32([]byte("aaaaaaaa"))
	assert.Equal(aa, float32(2.598459e+20))
}

func TestPutFloat32(t *testing.T) {
	assert := assert.New(t)
	a := []byte("  rrwwww")
	PutFloat32(a, float32(2.598459e+20))
	assert.Equal(a, []byte("aaaawwww"))
}

func TestFloat64(t *testing.T) {
	assert := assert.New(t)
	aa := Float64([]byte("aaaaaaaa"))
	assert.Equal(aa, float64(1.2217638442043777e+161))
}

func TestPutFloat64(t *testing.T) {
	assert := assert.New(t)
	a := []byte("  rrwwww")
	PutFloat64(a, float64(1.2217638442043777e+161))
	assert.Equal(a, []byte("aaaaaaaa"))
}
