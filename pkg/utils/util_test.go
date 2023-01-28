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

func Test_Min(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Min(1, 1), 1)
	assert.Equal(Min(1, 2), 1)
}

func Test_Max(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Max(1, 1), 1)
	assert.Equal(Max(1, 2), 2)
}

func Test_MinFloat64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(MinFloat64(1.2, 5.1), 1.2)
	assert.Equal(MinFloat64(1.2, 1.20), 1.2)
}

func Test_MaxFloat64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(MaxFloat64(1.2, 5.1), 5.1)
	assert.Equal(MaxFloat64(1.2, 1.20), 1.2)
}
