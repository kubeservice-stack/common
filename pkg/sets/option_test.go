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

package sets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	assert := assert.New(t)

	s := New(1, 2, 3)
	s2 := New(3, 2, 1)

	assert.True(s.Equal(s2))

	s3 := New("1", "2", "")
	s4 := New("1", "2", "2")
	s5 := New("1", "2")
	assert.False(s4.Equal(s3))
	assert.True(s4.Equal(s5))
}
