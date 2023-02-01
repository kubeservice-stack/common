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

package mcpack_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kubeservice-stack/common/pkg/codec/mcpack"
)

func TestDecodeEncodeMap(t *testing.T) {
	assert := assert.New(t)
	a := map[string]string{"aa": "bb"}
	m, err := Marshal(a)
	assert.Nil(err)
	assert.Equal(m, []uint8([]byte{0x10, 0x0, 0xd, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x3, 0x3, 0x61, 0x61, 0x0, 0x62, 0x62, 0x0}))
	b := new(map[string]string)
	err = Unmarshal(m, &b)
	assert.Nil(err)
	assert.Equal(*b, a)
}

type Person struct {
	Name string
}

func TestDecodeEncodeStruct(t *testing.T) {
	assert := assert.New(t)
	data, err := Marshal(Person{Name: "dongjiang"})
	assert.NoError(err)
	assert.Equal(data, []uint8([]byte{0x10, 0x0, 0x16, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x5, 0xa, 0x4e, 0x61, 0x6d, 0x65, 0x0, 0x64, 0x6f, 0x6e, 0x67, 0x6a, 0x69, 0x61, 0x6e, 0x67, 0x0}))
	p := &Person{}
	err = Unmarshal(data, p)
	assert.NoError(err)
	assert.Equal("dongjiang", p.Name)

}
