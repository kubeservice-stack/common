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

package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLimiter(t *testing.T) {
	assert := assert.New(t)
	aa, has := GetLimiter("RATELIMITER")
	assert.True(has)
	assert.NotNil(aa)

	bb, has := GetLimiter("empty")
	assert.False(has)
	assert.Nil(bb)
}

func Test_GetDefaultLimiter(t *testing.T) {
	assert := assert.New(t)

	bb := GetDefaultLimiter()
	assert.NotNil(bb)
}

func Test_HasRegister(t *testing.T) {
	assert := assert.New(t)

	ok := HasRegister("RATELIMITER")
	assert.True(ok)

	ok = HasRegister("empty")
	assert.False(ok)
}
