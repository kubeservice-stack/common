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

type info struct {
	AAA []int8
	BBB []int16
	CCC []int32
	DDD []int64
	EEE []uint16
	EAA []uint8
	FFF []uint32
	GGG []uint64
	HHH []float32
	III []float64
}

type Person struct {
	Name string
	Str  string
	Age  uint8
	Info info
	Ptr  uintptr
	Num  uint
	P    *int
}

func TestDecodeEncodeStruct(t *testing.T) {
	assert := assert.New(t)
	data, err := Marshal(
		Person{
			Name: "dongjiangAA Kk^2@4234",
			Str:  "!$_&-  éè  ;∞¥₤€",
			Age:  18,
			Ptr:  uintptr(1),
			Num:  111,
			Info: info{
				AAA: []int8{1, -1, 0, -0},
				BBB: []int16{345, 11, -1},
				CCC: []int32{123, 3, -333},
				DDD: []int64{111223, 44343, -333334, 0},
				EEE: []uint16{11223, 44343, 0},
				EAA: []uint8{12, 33, 33},
				FFF: []uint32{1199223, 44343, 0},
				GGG: []uint64{11999999223, 44343, 0},
				HHH: []float32{11999999223.994, 44343.333, 0.000, -3.1415926},
				III: []float64{-11999999223.994, 44343.333, 0.000, -3.1415926},
			},
		})
	assert.NoError(err)
	p := &Person{}
	err = Unmarshal(data, p)
	assert.NoError(err)
	assert.Equal("dongjiangAA Kk^2@4234", p.Name)
	assert.Equal("!$_&-  éè  ;∞¥₤€", p.Str)
	assert.Equal(uint8(18), p.Age)
	assert.Equal(uintptr(1), p.Ptr)
}

func TestDecodeEncodeNil(t *testing.T) {
	assert := assert.New(t)
	data, err := Marshal(
		Person{
			Name: "dongjiang",
			P:    nil,
		})
	assert.NoError(err)
	p := &Person{}
	err = Unmarshal(data, p)
	assert.NoError(err)
	assert.Equal("dongjiang", p.Name)

	p1 := &Person{}
	err = Unmarshal(data, p1)
	assert.NoError(err)
	assert.Equal("dongjiang", p.Name)
	assert.Equal((*int)(nil), p.P)

}
