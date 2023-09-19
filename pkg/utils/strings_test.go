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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetBetweenStr(t *testing.T) {
	assert := assert.New(t)

	substr := GetBetweenStr("aaaadfddweeacrttsdf", "a", "r")
	assert.Equal(substr, "aaaadfddweeac", "is true")

	substr = GetBetweenStr("aaaadfddweeacrttsdf", "d", "d")
	assert.Equal(substr, "", "is true")

	substr = GetBetweenStr("aaaadfddweeacrttsdf", "f", "t")
	assert.Equal(substr, "fddweeacr", "is true")
}

func Test_SubStr(t *testing.T) {
	assert := assert.New(t)

	substr := Substr("aaaadfddweeacrttsdf", -1, 0)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", 0, 1)
	assert.Equal(substr, "a", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, 1)
	assert.Equal(substr, "f", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, 2)
	assert.Equal(substr, "f", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, 3)
	assert.Equal(substr, "df", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2000, 3)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, -3)
	assert.Equal(substr, "tts", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -2, -100)
	assert.Equal(substr, "aaaalfbkweeacrtts", "is true")

	substr = Substr("aaaalfbkweeacrttsdf", -1, -100)
	assert.Equal(substr, "aaaalfbkweeacrttsd", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 1000, 11111)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 0)
	assert.Equal(substr, "", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 10000)
	assert.Equal(substr, "dfddweeacrttsdf", "is true")

	substr = Substr("aaaadfddweeacrttsdf", 4, 6)
	assert.Equal(substr, "dfddwe", "is true")
}

func Test_Strings2Bytes(t *testing.T) {
	assert := assert.New(t)
	b := Strings2Bytes("adfbd中国 adsf")
	assert.Equal(b, []byte{0x61, 0x64, 0x66, 0x62, 0x64, 0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd, 0x20, 0x61, 0x64, 0x73, 0x66})

	b = Strings2Bytes("")
	assert.Equal(b, []byte(nil))

	s := new(string)
	*s = "aaa"
	b = Strings2Bytes(*s)
	assert.Equal(b, []byte{0x61, 0x61, 0x61})

	*s = ""
	b = Strings2Bytes(*s)
	assert.Equal(b, []byte(nil))
}

func Test_Bytes2Strings(t *testing.T) {
	assert := assert.New(t)
	s := Bytes2Strings([]byte{0x61, 0x64, 0x66, 0x62, 0x64, 0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd, 0x20, 0x61, 0x64, 0x73, 0x66})
	assert.Equal(s, "adfbd中国 adsf")

	b := new([]byte)
	*b = nil
	s = Bytes2Strings(*b)
	assert.Equal(s, "")

	*b = []byte{}
	s = Bytes2Strings(*b)
	assert.Equal(s, "")

	*b = []byte("!$_&-  éè  ;∞¥₤€")
	s = Bytes2Strings(*b)
	assert.Equal(s, "!$_&-  éè  ;∞¥₤€")
}
