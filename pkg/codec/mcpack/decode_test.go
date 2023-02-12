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
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type unmarshalTest struct {
	in        []byte
	ptr       interface{}
	out       interface{}
	equalFunc func(interface{}, interface{}) bool
}

type obj struct {
	Foo string `json:"foo"`
}

type ping struct {
	Data string
}

type UV struct {
	F1 *UU    `json:"F1"`
	F2 int32  `json:"F2"`
	F3 Number `json:"F3"`
}

type UU struct {
	Data []byte
}

func (u *UU) UnmarshalMCPACK(b []byte) error {
	u.Data = b
	return nil
}

type Ptr struct {
	I8   *int8    `json:"i8"`
	I16  *int16   `json:"i16"`
	I32  *int32   `json:"i32"`
	I64  *int64   `json:"i64"`
	UI8  *uint8   `json:"ui8"`
	UI16 *uint16  `json:"ui16"`
	UI32 *uint32  `json:"ui32"`
	UI64 *uint64  `json:"ui64"`
	F32  *float32 `json:"f32"`
	F64  *float64 `json:"f64"`
	Arr  []*int   `json:"arr"`
}

var unmarshalTests = []unmarshalTest{
	{
		in: []byte{MCPACKV2_OBJECT, 0, 20, 0, 0, 0, 1, 0, 0, 0, MCPACKV2_STRING, 6, 4, 0, 0, 0, 'a', 'l', 'p', 'h', 'a', 0, 'a', '-', 'z', 0},

		ptr: &UU{},
		out: &UU{Data: []byte{MCPACKV2_OBJECT, 0, 20, 0, 0, 0, 1, 0, 0, 0, MCPACKV2_STRING, 6, 4, 0, 0, 0, 'a', 'l', 'p', 'h', 'a', 0, 'a', '-', 'z', 0}},
		equalFunc: func(l, r interface{}) bool {
			var ll, rr *UU = l.(*UU), r.(*UU)
			return bytes.Equal(ll.Data, rr.Data)
		},
	},
	{
		in:  []byte{MCPACKV2_STRING, 0, 4, 0, 0, 0, 'f', 'o', 'o', 0},
		ptr: new(string),
		out: "foo",
	},
	{
		in:  []byte{MCPACKV2_INT32, 0, 4, 0, 0, 0},
		ptr: new(int32),
		out: int32(4),
	},
	{
		in:  []byte{MCPACKV2_INT64, 0, 4, 0, 0, 0, 0, 0, 0, 0},
		ptr: new(int64),
		out: int64(4),
	},
	{
		in:  []byte{MCPACKV2_BOOL, 0, 1},
		ptr: new(bool),
		out: bool(true),
	},
	{
		in:  []byte{MCPACKV2_BOOL, 0, 0},
		ptr: new(bool),
		out: bool(false),
	},
	{
		in: []byte{MCPACKV2_OBJECT, 0, 0, 0, 0, 0, 1, 0, 0, 0,
			MCPACKV2_STRING, 4, 4, 0, 0, 0, 'f', 'o', 'o', 0, 'b', 'a', 'r', 0},
		ptr: new(obj),
		out: obj{Foo: "bar"},
	},
	{
		in: []byte{MCPACKV2_OBJECT, 0, 0, 0, 0, 0, 1, 0, 0, 0,
			MCPACKV2_SHORT_STRING, 4, 4, 'f', 'o', 'o', 0, 'b', 'a', 'r', 0},
		ptr: new(obj),
		out: obj{Foo: "bar"},
	},
	{
		in: []byte{MCPACKV2_ARRAY, 0, 0, 0, 0, 0, 1, 0, 0, 0,
			MCPACKV2_STRING, 0, 4, 0, 0, 0, 'f', 'o', 'o', 0},
		ptr: new([]string),
		out: []string{"foo"},
	},
	{
		in:  []byte{MCPACKV2_OBJECT, 0, 17, 0, 0, 0, 1, 0, 0, 0, 208, 5, 5, 'D', 'a', 't', 'a', 0, 'p', 'i', 'n', 'g', 0},
		ptr: new(ping),
		out: ping{Data: "ping"},
	},
}

func TestUnmarshal(t *testing.T) {
	assert := assert.New(t)
	for _, tt := range unmarshalTests {
		err := Unmarshal(tt.in, tt.ptr)
		assert.Nil(err)

		if tt.equalFunc != nil {
			assert.False(!tt.equalFunc(tt.ptr, tt.out))
		} else {
			assert.False(!reflect.DeepEqual(reflect.ValueOf(tt.ptr).Elem().Interface(), tt.out))
		}
	}
}

func TestUnmarshalObeject(t *testing.T) {
	assert := assert.New(t)
	obj := new(T)
	err := Unmarshal([]byte{MCPACKV2_OBJECT, 0, 28, 0, 0, 0,
		3, 0, 0, 0,
		MCPACKV2_BOOL, 2, 'A', 0, 1,
		MCPACKV2_SHORT_STRING, 2, 2, 'X', 0, 'x', 0,
		MCPACKV2_INT64, 2, 'Y', 0, 1, 0, 0, 0, 0, 0, 0, 0}, obj)
	assert.Nil(err)
	assert.Equal(obj, &T{A: true, X: "x", Y: 1, Z: 0})

	objv := new(V)
	err = Unmarshal([]byte{MCPACKV2_OBJECT, 0, 52, 0, 0, 0, 3, 0, 0, 0,
		MCPACKV2_OBJECT, 3, 17, 0, 0, 0, 'F', '1', 0, 1, 0, 0, 0, MCPACKV2_SHORT_STRING, 6, 4, 'a', 'l', 'p', 'h', 'a', 0, 'a', '-', 'z', 0,
		MCPACKV2_INT32, 3, 'F', '2', 0, 1, 0, 0, 0,
		MCPACKV2_INT64, 3, 'F', '3', 0, 1, 0, 0, 0, 0, 0, 0, 0}, objv)
	assert.Nil(err)
	assert.Equal(objv, &V{F1: map[string]interface{}{"alpha": "a-z"}, F2: 1, F3: Number(1)})

	objy := new(Y)
	err = Unmarshal([]byte{MCPACKV2_OBJECT, 0, 20, 0, 0, 0,
		1, 0, 0, 0,
		MCPACKV2_ARRAY, 6, 4, 0, 0, 0, 'E', 'm', 'p', 't', 'y', 0, 0, 0, 0, 0,
	}, objy)
	assert.Nil(err)
	assert.Equal(objy, &Y{Empty: []string{}})

	a := []interface{}{&ping{Data: "a"}, nil}
	te, err := Marshal(a)
	assert.Nil(err)
	assert.Equal(te, []byte{0x20, 0x0, 0x1b, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0xe, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x0, 0x61, 0x0, 0x0})

	var b []interface{}
	err = Unmarshal([]byte{0x20, 0x0, 0x1b, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0xe, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x0, 0x61, 0x0, 0x0}, b)
	assert.NotNil(err)
	assert.Equal(err.Error(), "mcpack: Unmarshal(non-pointer []interface {})")

	var c *int
	err = Unmarshal([]byte{0x20, 0x0, 0x1b, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0xe, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x0, 0x61, 0x0, 0x0}, c)
	assert.NotNil(err)
	assert.Equal(err.Error(), "mcpack: Unmarshal(nil *int)")

	var d []*interface{}
	err = Unmarshal([]byte{0x20, 0x0, 0x1b, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0xe, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xd0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x0, 0x61, 0x0, 0x0}, d)
	assert.NotNil(err)
	assert.Equal(err.Error(), "mcpack: Unmarshal(non-pointer []*interface {})")
}

func TestUV(t *testing.T) {
	assert := assert.New(t)
	a := UV{F1: &UU{Data: []byte("aa")}, F2: int32(33), F3: Number(11)}
	te, err := Marshal(a)
	assert.Nil(err)
	assert.Equal(te, []byte{0x10, 0x0, 0x31, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x10, 0x3, 0xe, 0x0, 0x0, 0x0, 0x46, 0x31, 0x0, 0x1, 0x0, 0x0, 0x0, 0xe0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x61, 0x14, 0x3, 0x46, 0x32, 0x0, 0x21, 0x0, 0x0, 0x0, 0x18, 0x3, 0x46, 0x33, 0x0, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})

	b := new(UV)
	err = Unmarshal([]byte{0x10, 0x0, 0x31, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x10, 0x3, 0xe, 0x0, 0x0, 0x0, 0x46, 0x31, 0x0, 0x1, 0x0, 0x0, 0x0, 0xe0, 0x5, 0x2, 0x44, 0x61, 0x74, 0x61, 0x0, 0x61, 0x61, 0x14, 0x3, 0x46, 0x32, 0x0, 0x21, 0x0, 0x0, 0x0, 0x18, 0x3, 0x46, 0x33, 0x0, 0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, b)
	assert.Nil(err)
	assert.Equal(b.F2, int32(33))
	assert.Equal(b.F3, Number(11))
}

func TestPrt(t *testing.T) {
	assert := assert.New(t)
	value := new(int8)
	*value = 1
	a := Ptr{
		I8: value,
	}
	te, err := Marshal(a)
	assert.Nil(err)
	assert.Equal(te, []byte{0x10, 0x0, 0x5a, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x11, 0x3, 0x69, 0x38, 0x0, 0x1, 0x61, 0x4, 0x69, 0x31, 0x36, 0x0, 0x0, 0x61, 0x4, 0x69, 0x33, 0x32, 0x0, 0x0, 0x61, 0x4, 0x69, 0x36, 0x34, 0x0, 0x0, 0x61, 0x4, 0x75, 0x69, 0x38, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x31, 0x36, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x33, 0x32, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x36, 0x34, 0x0, 0x0, 0x61, 0x4, 0x66, 0x33, 0x32, 0x0, 0x0, 0x61, 0x4, 0x66, 0x36, 0x34, 0x0, 0x0, 0x20, 0x4, 0x4, 0x0, 0x0, 0x0, 0x61, 0x72, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0})

	b := new(Ptr)
	err = Unmarshal([]byte{0x10, 0x0, 0x5a, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x11, 0x3, 0x69, 0x38, 0x0, 0x1, 0x61, 0x4, 0x69, 0x31, 0x36, 0x0, 0x0, 0x61, 0x4, 0x69, 0x33, 0x32, 0x0, 0x0, 0x61, 0x4, 0x69, 0x36, 0x34, 0x0, 0x0, 0x61, 0x4, 0x75, 0x69, 0x38, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x31, 0x36, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x33, 0x32, 0x0, 0x0, 0x61, 0x5, 0x75, 0x69, 0x36, 0x34, 0x0, 0x0, 0x61, 0x4, 0x66, 0x33, 0x32, 0x0, 0x0, 0x61, 0x4, 0x66, 0x36, 0x34, 0x0, 0x0, 0x20, 0x4, 0x4, 0x0, 0x0, 0x0, 0x61, 0x72, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0}, b)
	assert.Nil(err)
	assert.Equal(*b.I8, int8(1))
	assert.Equal(b.I16, (*int16)(nil))
	assert.Equal(b.Arr, []*int{})
}

func TestFloat32Decodex(t *testing.T) {
	assert := assert.New(t)
	va := new(int)
	*va = 1
	str := new(string)
	*str = "aadd"
	str1 := new(string)
	*str1 = "Users/dongjiang/Documentsdfasdntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github"
	a := &V{
		F1: map[string]interface{}{
			"ui64":  uint64(0xFFFFFFFFFFFFFFFF),
			"ui32":  uint32(0xFFFFFFFF),
			"bys":   bytes.Runes([]byte("dasdf")),
			"alpha": "a-z",
			"a":     1,
			"is":    true,
			"str":   str,
			"tttt":  str1,
			"dd":    float64(1.11),
			"ff":    float32(1.11),
			"b":     va,
			"c":     reflect.ValueOf(va),
			"d":     map[string]interface{}{"aa": "bb"},
			"e":     []interface{}{"aa", 1, va},
			"f":     map[string]float64{"a": float64(-45.2231)},
		},
		F2: 1,
		F3: Number(1),
	}
	te, err := Marshal(a)
	assert.Nil(err)
	b := new(V)
	err = Unmarshal(te, b)
	assert.Nil(err)
	assert.Equal(b, &V{
		F1: map[string]interface{}{
			"a":     int64(1),
			"bys":   []interface{}{int32(100), int32(97), int32(115), int32(100), int32(102)},
			"ui64":  uint64(18446744073709551615),
			"ui32":  uint32(4294967295),
			"alpha": "a-z",
			"str":   "aadd",
			"tttt":  *str1,
			"is":    true,
			"dd":    float64(1.11),
			"ff":    float32(1.11),
			"b":     int64(1),
			"c":     map[string]interface{}{},
			"d":     map[string]interface{}{"aa": "bb"},
			"e":     []interface{}{"aa", int64(1), int64(1)},
			"f":     map[string]interface{}{"a": -45.2231},
		},
		F2: 1,
		F3: 1,
	})

	c := make([]byte, 5)
	err = Unmarshal([]byte{MCPACKV2_BINARY, 0, 4, 0, 0, 0, 'f', 'o', 'o', 0}, &c)
	assert.Nil(err)
	assert.Equal(c, []byte{'f', 'o', 'o', 0})

	d := make([]byte, 10)
	err = Unmarshal([]byte{MCPACKV2_SHORT_BINARY, 4, 4, 'f', 'o', 'o', 0, 'b', 'a', 'r', 0}, &d)
	assert.Nil(err)
	assert.Equal(d, []byte{'b', 'a', 'r', 0})

	e := make([]byte, 10)
	err = Unmarshal([]byte{MCPACKV2_NULL, 0, 0}, &e)
	assert.Nil(err)
	assert.Equal(e, []byte(nil))
}

func TestDecodeState(t *testing.T) {
	assert := assert.New(t)
	d := decodeState{
		data: []byte("111asddfsdffdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddasdfasdfsadfsdf111asddfsdffdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddasdfasdfsadfsdf"),
		off:  1,
	}

	aa := d.shortBinaryInterface()
	assert.Equal(aa, []byte{0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x61, 0x73, 0x64, 0x66, 0x61, 0x73, 0x64, 0x66, 0x73, 0x61, 0x64, 0x66, 0x73, 0x64, 0x66, 0x31, 0x31, 0x31, 0x61, 0x73, 0x64, 0x64, 0x66, 0x73, 0x64, 0x66, 0x66, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64, 0x64})

	bb := d.nullInterface()
	assert.Equal(bb, nil)
}
