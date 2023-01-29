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

package errno

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	t1 := New(0, "ok")
	t2 := New(0, "ok")
	t3 := New(-1, "message")

	assert.Equal(t1, t2)

	assert.Equal(t1.status, 0)
	assert.Equal(t1.message, "ok")

	assert.Equal(t1.Status(), 0)
	assert.Equal(t1.Message(), "ok")

	assert.Equal(t3.status, -1)
	assert.Equal(t3.message, "message")
}

func TestNewCode(t *testing.T) {
	assert := assert.New(t)
	t1 := NewCode(0)
	t2 := NewCode(0)
	t3 := NewCode(-1)

	assert.Equal(t1, t2)

	assert.Equal(t1.status, 0)
	assert.Equal(t1.message, "")

	assert.Equal(t2.status, 0)
	assert.Equal(t2.message, "")

	assert.Equal(t3.status, -1)
	assert.Equal(t3.message, "")
}

func TestErrno_Error(t *testing.T) {
	assert := assert.New(t)
	err := New(-1, "System Error")

	assert.Equal(err.Error(), "Error - errno: -1, errmsg: System Error")
	assert.Equal(fmt.Sprintf("%v", err), "Error - errno: -1, errmsg: System Error")
}
