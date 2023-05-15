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

package schedule

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var functionNameA = func() bool {
	return true
}

func functionNameC() bool {
	return true
}

func Test_GetFunctionName(t *testing.T) {
	assert := assert.New(t)
	functionNameB := func() bool {
		return true
	}
	aa := getFunctionName(functionNameA)
	assert.Contains(aa, ".func1")
	bb := getFunctionName(functionNameB)
	assert.Contains(bb, "Test_GetFunctionName.func1")
	cc := getFunctionName(functionNameC)
	assert.Contains(cc, "functionNameC")
}

func Test_GetFunctionKey(t *testing.T) {
	assert := assert.New(t)
	aa := getFunctionKey("ddd")
	assert.Equal(aa, "730f75dafd73e047b86acb2dbd74e75dcb93272fa084a9082848f2341aa1abb6")
	bb := getFunctionKey("")
	assert.Equal(bb, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func Test_CallTaskFuncWithParams(t *testing.T) {
	assert := assert.New(t)

	aa, err := callTaskFuncWithParams(functionNameA, nil)
	assert.Nil(err)
	assert.NotNil(aa)

	bb, err := callTaskFuncWithParams(functionNameA, []interface{}{1, 1})
	assert.Equal(err, ErrParamsNotAdapted)
	assert.Equal(bb, []reflect.Value([]reflect.Value(nil)))

	cc, err := callTaskFuncWithParams(functionNameC, nil)
	assert.Nil(err)
	assert.NotNil(cc)

	dd, err := callTaskFuncWithParams(functionNameC, nil)
	assert.Nil(err)
	assert.NotNil(dd)

	functionNameB := func() bool {
		return true
	}

	ee, err := callTaskFuncWithParams(functionNameB, nil)
	assert.Nil(err)
	assert.NotNil(ee)

	functionNot := strings.Fields
	ff, err := callTaskFuncWithParams(functionNot, []interface{}{"dfadf"})
	assert.Nil(err)
	assert.NotNil(ff)
}

func Test_formatTime(t *testing.T) {
	assert := assert.New(t)
	h, m, s, err := formatTime("10:12")
	assert.Nil(err)
	assert.Equal(h, 10)
	assert.Equal(m, 12)
	assert.Equal(s, 0)

	h, m, s, err = formatTime("23:59:59")
	assert.Nil(err)
	assert.Equal(h, 23)
	assert.Equal(m, 59)
	assert.Equal(s, 59)
	h, m, s, err = formatTime("0:0:0")
	assert.Nil(err)
	assert.Equal(h, 0)
	assert.Equal(m, 0)
	assert.Equal(s, 0)

	h, m, s, err = formatTime("24:59:59")
	assert.NotNil(err)
	h, m, s, err = formatTime("23:59:59:11")
	assert.NotNil(err)
	h, m, s, err = formatTime("er")
	assert.NotNil(err)
	h, m, s, err = formatTime("11")
	assert.NotNil(err)
	h, m, s, err = formatTime("aa:59:59")
	assert.NotNil(err)
	h, m, s, err = formatTime("1:a:59")
	assert.NotNil(err)
	h, m, s, err = formatTime("1:11:a")
	assert.NotNil(err)

	h, m, s, err = formatTime("23:59:59.323")
	assert.NotNil(err)
}
