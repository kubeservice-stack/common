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
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type marshalTest struct {
	in  interface{}
	out []byte
}

type T struct {
	A bool
	X string
	Y int
	Z int `json:"-"`
}

type U struct {
	Alphabet string `json:"alpha"`
}

type V struct {
	F1 interface{}
	F2 int32
	F3 Number
}

type Number int

type W struct {
	S string
	V int16
}

type X struct {
	Beta map[string]string
	Deta [2]int16
}

type Y struct {
	Empty []string
}

type E struct {
	Beta map[string]string
}

func (e *E) Error() string {
	return fmt.Sprintf("len(key) exceeds %d", MCPACKV2_KEY_MAX_LEN)
}

var longVItem = [299]byte{250: '1', 251: '8', 252: '2', 253: '2', 297: 'S', 298: 'V'}

var marshalTests = []marshalTest{
	{
		in: &T{A: true, X: "x", Y: 1, Z: 2},
		out: []byte{MCPACKV2_OBJECT, 0, 28, 0, 0, 0,
			3, 0, 0, 0,
			MCPACKV2_BOOL, 2, 'A', 0, 1,
			MCPACKV2_SHORT_STRING, 2, 2, 'X', 0, 'x', 0,
			MCPACKV2_INT64, 2, 'Y', 0, 1, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		in: &U{Alphabet: "a-z"},
		out: []byte{MCPACKV2_OBJECT, 0, 17, 0, 0, 0,
			1, 0, 0, 0,
			MCPACKV2_SHORT_STRING, 6, 4, 'a', 'l', 'p', 'h', 'a', 0, 'a', '-', 'z', 0},
	},
	{
		in: &V{F1: &U{Alphabet: "a-z"}, F2: 1, F3: Number(1)},
		out: []byte{MCPACKV2_OBJECT, 0, 52, 0, 0, 0,
			3, 0, 0, 0,
			MCPACKV2_OBJECT, 3, 17, 0, 0, 0, 'F', '1', 0, 1, 0, 0, 0, MCPACKV2_SHORT_STRING, 6, 4, 'a', 'l', 'p', 'h', 'a', 0, 'a', '-', 'z', 0,
			MCPACKV2_INT32, 3, 'F', '2', 0, 1, 0, 0, 0,
			MCPACKV2_INT64, 3, 'F', '3', 0, 1, 0, 0, 0, 0, 0, 0, 0},
	},
	getTestslongVItemW(),
	getTestsKeyBoundaryX(),
	{
		in: &Y{},
		out: []byte{MCPACKV2_OBJECT, 0, 20, 0, 0, 0,
			1, 0, 0, 0,
			MCPACKV2_ARRAY, 6, 4, 0, 0, 0, 'E', 'm', 'p', 't', 'y', 0, 0, 0, 0, 0,
		},
	},
	getTestsKeyTooLongE(),
}

func getTestslongVItemW() marshalTest {
	out := []byte{MCPACKV2_OBJECT, 0, 0x40, 0x1, 0, 0,
		2, 0, 0, 0,
		MCPACKV2_STRING, 2, 0x2c, 0x1, 0, 0, 'S', 0,
	}
	out = append(out, longVItem[:]...)
	out = append(out, 0, MCPACKV2_INT32, 2, 'V', 0, 1, 0, 0, 0)
	return marshalTest{
		in:  &W{S: string(longVItem[:]), V: 1},
		out: out,
	}
}

func getTestsKeyBoundaryX() marshalTest {
	beta := make(map[string]string)
	key := string(longVItem[:254])
	beta[key] = "SV"

	deta := [2]int16{-18, -22}
	in := &X{Beta: beta, Deta: deta}

	out := []byte{MCPACKV2_OBJECT, 0, 0x33, 0x1, 0, 0, //header
		2, 0, 0, 0,
		MCPACKV2_OBJECT, 5, 0x9, 0x1, 0, 0, //map
		'B', 'e', 't', 'a', 0, //key: Beta | 0x0
		1, 0, 0, 0, //count: 1
		MCPACKV2_SHORT_STRING, 0xff, 3}
	out = append(out, longVItem[:254]...)
	out = append(out, 0, 'S', 'V', 0)
	//Deta
	out = append(out, MCPACKV2_ARRAY, 5, 0x10, 0, 0, 0,
		'D', 'e', 't', 'a', 0,
		2, 0, 0, 0,
		MCPACKV2_INT32, 0, 0xee, 0xff, 0xff, 0xff, MCPACKV2_INT32, 0, 0xea, 0xff, 0xff, 0xff)

	return marshalTest{
		in:  in,
		out: out,
	}
}

func getTestsKeyTooLongE() marshalTest {
	beta := make(map[string]string)
	key := string(longVItem[:255])
	beta[key] = "SV"
	in := &E{Beta: beta}

	return marshalTest{
		in:  in,
		out: nil, // an nil out impiles an error
	}
}

func TestMarshal(t *testing.T) {
	assert := assert.New(t)
	for _, tt := range marshalTests {
		b, err := Marshal(tt.in)
		// an error expected
		if tt.out == nil {
			if err == nil {
				assert.Nil(err)
			} else if tt.in.(error).Error() != err.Error() {
				assert.Equal(tt.in.(error).Error(), err.Error())
			}
			continue
		}

		// failed
		if err != nil {
			assert.Nil(err)
		}
		if !bytes.Equal(tt.out, b) {
			assert.Equal(tt.out, b)
		}

	}
}

func TestDominantField(t *testing.T) {
	assert := assert.New(t)
	aa, ok := dominantField([]field{
		field{
			name:      "aa",
			nameBytes: []byte("aa"),
		},
	})
	assert.True(ok)
	assert.Equal(aa, field{
		name:      "aa",
		nameBytes: []byte("aa"),
	})

	bb, ok := dominantField([]field{
		field{
			name:      "aa",
			nameBytes: []byte("aa"),
		},
		field{
			name:      "bb",
			nameBytes: []byte("bb"),
		},
	})
	assert.False(ok)
	assert.Equal(bb, field{})

	cc, ok := dominantField([]field{
		field{
			index: []int{10, 11, 12},
		},
	})
	assert.True(ok)
	assert.Equal(cc, field{
		index: []int{10, 11, 12},
	})
}

func TestIsEmptyValue(t *testing.T) {
	assert := assert.New(t)
	var s []string
	s = []string{"aa", "bb"}
	ok := isEmptyValue(reflect.ValueOf(s))
	assert.False(ok)

	var a int
	ok = isEmptyValue(reflect.ValueOf(a))
	assert.True(ok)

	var b map[string]string
	ok = isEmptyValue(reflect.ValueOf(b))
	assert.True(ok)

	var c bool
	ok = isEmptyValue(reflect.ValueOf(c))
	assert.True(ok)

	var d uint32
	ok = isEmptyValue(reflect.ValueOf(d))
	assert.True(ok)

	var e float32
	ok = isEmptyValue(reflect.ValueOf(e))
	assert.True(ok)

	var f interface{}
	ok = isEmptyValue(reflect.ValueOf(f))
	assert.False(ok)

	var g T
	ok = isEmptyValue(reflect.ValueOf(g))
	assert.False(ok)

	var h []int
	ok = isEmptyValue(reflect.ValueOf(h))
	assert.True(ok)

	var i rune
	ok = isEmptyValue(reflect.ValueOf(i))
	assert.True(ok)
}

func TestBinaryEncoder(t *testing.T) {
	assert := assert.New(t)
	e := encodeState{
		data: []byte("aaaaaaaä"),
		off:  1,
	}
	var s []byte
	s = []byte("adfasdf")
	binaryEncoder(&e, "dfasdvsdghhd", reflect.ValueOf(s))
	assert.Equal(e.off, 24)

	var ss []byte
	ss = []byte("adfassdfasdfsdfsdfdfgsadfasdfsdfdf/opt/homebrew/Cellar/go/1.19.4/libexec/bin/go test -v [/Users/dongjiang/Documentsdfasdntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasdfasdf")
	binaryEncoder(&e, "dff", reflect.ValueOf(ss))
	assert.Equal(e.off, 1192)
}
