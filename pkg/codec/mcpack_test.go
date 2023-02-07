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

package codec_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/kubeservice-stack/common/pkg/codec"
	"github.com/stretchr/testify/assert"
)

type TV struct {
	F interface{} `json:"f"`
}

func TestMCPack(t *testing.T) {
	assert := assert.New(t)
	va := new(int)
	*va = 1
	str := new(string)
	*str = "dongjiang"
	str1 := new(string)
	*str1 = "long string users/dongjiang/Documentsdfasdntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github.com/kubeservice-stack/common/pkg/codec/mcpack]asdfasntsdfasdfsdgs/go/src/github"
	a := &TV{
		F: map[string]interface{}{
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
	}
	te, err := codec.PluginInstance(codec.MCPACK).Marshal(a)
	assert.Nil(err)
	b := new(TV)
	err = codec.PluginInstance(codec.MCPACK).Unmarshal(te, b)
	assert.Nil(err)
	assert.Equal(b, &TV{
		F: map[string]interface{}{
			"a":     int64(1),
			"bys":   []interface{}{int32(100), int32(97), int32(115), int32(100), int32(102)},
			"ui64":  uint64(18446744073709551615),
			"ui32":  uint32(4294967295),
			"alpha": "a-z",
			"str":   "dongjiang",
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
	})
}
