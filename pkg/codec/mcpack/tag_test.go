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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatTag(t *testing.T) {
	assert := assert.New(t)
	aa := tagOptions("aaaa,aaaa,,bbb,")
	bb := tagOptions("")
	ok := bb.Contains("")
	assert.False(ok)

	ok = bb.Contains(string([]byte{}))
	assert.False(ok)

	ok = aa.Contains("aa")
	assert.False(ok)

	ok = aa.Contains("aaaa")
	assert.True(ok)

	ok = aa.Contains("")
	assert.True(ok)

	a, b := parseTag("aa,bb,,c")
	assert.Equal(a, "aa")
	assert.Equal(b, tagOptions("bb,,c"))
}
