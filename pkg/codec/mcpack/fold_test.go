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

func TestFoldFunc(t *testing.T) {
	assert := assert.New(t)
	var aaa func(s, t []byte) bool
	var bbb func(s, t []byte) bool
	aaa = foldFunc([]byte("# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.com/kubeservice-stack/common/pkg/codec/mcpack.test]"))
	assert.False(aaa([]byte("aa"), []byte("bb")))
	assert.False(aaa([]byte(""), []byte("bb")))
	assert.True(aaa([]byte("bb"), []byte("bb")))
	assert.True(aaa([]byte("Abb"), []byte("Abb")))

	bbb = foldFunc([]byte(""))
	assert.False(bbb([]byte("aa"), []byte("bb")))
	assert.False(bbb([]byte(""), []byte("bb")))
	assert.True(bbb([]byte("bb"), []byte("bb")))
	assert.True(bbb([]byte("Abb"), []byte("Abb")))
}

func TestEqualFoldRight(t *testing.T) {
	assert := assert.New(t)
	assert.True(equalFoldRight([]byte("a"), []byte("a")))
	assert.True(equalFoldRight([]byte(""), []byte("")))
	assert.False(equalFoldRight([]byte(""), []byte("11")))
	assert.True(equalFoldRight([]byte("A"), []byte("a")))
	assert.False(equalFoldRight([]byte("b"), []byte("a")))
	assert.True(equalFoldRight([]byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA")))
	assert.False(equalFoldRight([]byte("!# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA")))
	assert.True(equalFoldRight([]byte("∞¥₤€"), []byte("∞¥₤€")))
	assert.False(equalFoldRight([]byte("∞¥₤€"), []byte("")))
}

func TestAsciiEqualFold(t *testing.T) {
	assert := assert.New(t)
	assert.True(asciiEqualFold([]byte("a"), []byte("a")))
	assert.True(asciiEqualFold([]byte("A"), []byte("a")))
	assert.False(asciiEqualFold([]byte("b"), []byte("a")))
	assert.True(asciiEqualFold([]byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA")))
	assert.False(asciiEqualFold([]byte("!# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA")))
}

func TestSimpleLetterEqualFold(t *testing.T) {
	assert := assert.New(t)
	assert.True(simpleLetterEqualFold([]byte("a"), []byte("a")))
	assert.True(simpleLetterEqualFold([]byte("A"), []byte("a")))
	assert.False(simpleLetterEqualFold([]byte("b"), []byte("a")))
	assert.True(simpleLetterEqualFold([]byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpacKS [github.coA")))
	assert.False(simpleLetterEqualFold([]byte("!# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA"), []byte("# github.com/kubeservice-stack/common/pkg/codec/mcpack [github.coA")))
}
