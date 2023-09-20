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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_ToStringDict(t *testing.T) {
	assert := assert.New(t)
	aastr := `{"err_no":0,"err_msg":"","data":[{"phone":"18500844425"},{"phone":"11000003616"},{"phone":"11000001843"}]}`

	var data map[string]interface{}

	err := json.Unmarshal([]byte(aastr), &data)
	assert.Nil(err, "is not err")

	da, ok := data["data"]
	assert.True(ok)

	plist, err := ToSlice(da)
	assert.Nil(err)

	aa, err := ToStringDict(plist, "phone")
	assert.Nil(err, "is not err")
	assert.NotNil(aa, "is not err")
}

func TestUtil_InSlice(t *testing.T) {
	assert := assert.New(t)

	assert.True(InSlice("aaa", []string{"aaaa", "aaa"}), "解析文件正常，不符合需求！")

	assert.False(InSlice("b", []string{"aaaa", "aaa"}), "解析文件正常，不符合需求！")

	assert.False(InSlice("aaa", []string{}), "解析文件正常，不符合需求！")

}

func Test_ReplayStr(t *testing.T) {
	assert := assert.New(t)
	r := ReplayStr(2, 2)
	assert.Equal(r, "02")

	r = ReplayStr(3242354456777986, 9)
	assert.Equal(r, "3242354456777986")
}

func TestUtil_InSliceIfaceToLower(t *testing.T) {
	assert := assert.New(t)

	tmp := make([]interface{}, 10)
	tmp = append(tmp, "aaa")

	ret, err := InSliceIfaceToLower("aaa", tmp)
	assert.Nil(err, "is not err")
	assert.True(ret, "解析文件正常，不符合需求！")

	ret, err = InSliceIfaceToLower("aAa", tmp)
	assert.Nil(err, "is not err")
	assert.True(ret, "解析文件正常，不符合需求！")

	ret, err = InSliceIfaceToLower("b", tmp)
	assert.Nil(err, "is not err")
	assert.False(ret, "解析文件正常，不符合需求！")

	ret, err = InSliceIfaceToLower("b", nil)
	assert.NotNil(err, "is err")
	assert.False(ret, "解析文件正常，不符合需求！")

	ret, err = InSliceIfaceToLower("b", []string{"aa", "bb"})
	assert.Nil(err, "is err")
	assert.False(ret, "解析文件正常，不符合需求！")

	ret, err = InSliceIfaceToLower("b", map[string]string{"aaaa": "aaa", "bbbb": "bbb"})
	assert.NotNil(err, "is err")
	assert.False(ret, "解析文件正常，不符合需求！")
}

func Test_SliceFilter(t *testing.T) {
	assert := assert.New(t)
	aaa := SliceFilter([]interface{}{"aa", 1, nil}, func(v interface{}) bool {
		if v == nil {
			return false
		} else {
			return true
		}
	})

	assert.Equal(aaa, []interface{}{"aa", 1})
}

func Test_SliceSum(t *testing.T) {
	assert := assert.New(t)
	sum := SliceSum([]int64{12, 34, 5, 224})
	assert.Equal(sum, int64(275))
}

func Test_SliceRange(t *testing.T) {
	assert := assert.New(t)
	aa := SliceRange(1, 100, 2)
	assert.Equal(aa, []int64([]int64{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49, 51, 53, 55, 57, 59, 61, 63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 85, 87, 89, 91, 93, 95, 97, 99}))

	bb := SliceRange(1, 100, 101)
	assert.Equal(bb, []int64{1})

	cc := SliceRange(1000, 100, 2)
	assert.Equal(cc, []int64(nil))
}

func Test_SliceRand(t *testing.T) {
	assert := assert.New(t)
	a := SliceRand([]interface{}{"aa"})
	assert.Equal(a, "aa")

	a = SliceRand([]interface{}{"aa", "bb"})
	assert.Contains("aabb", a)
}

func Test_InterfacesToStrings(t *testing.T) {
	assert := assert.New(t)
	aa := InterfacesToStrings([]interface{}{"aa", "123", "3dv3"})
	assert.Equal([]string{"aa", "123", "3dv3"}, aa)
}

func Test_SliceReduce(t *testing.T) {
	assert := assert.New(t)
	aa := SliceReduce([]interface{}{"aa", "123", "3dv3"}, func(v interface{}) interface{} {
		if v == "aa" {
			return false
		} else {
			return true
		}
	})
	assert.Equal([]interface{}([]interface{}{false, true, true}), aa)
}

func Test_SliceMerge(t *testing.T) {
	assert := assert.New(t)
	aa := SliceMerge([]interface{}{"aa", "123", "3dv3"}, []interface{}{"aa", "123", "3dv3"})
	assert.Equal([]interface{}{"aa", "123", "3dv3", "aa", "123", "3dv3"}, aa)
}
