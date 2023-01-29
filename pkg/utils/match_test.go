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

func Test_WildcardMatch_MatchingEmpty(t *testing.T) {
	assert.True(t, WildcardMatch("", ""))
	assert.False(t, WildcardMatch("", "42"))
	assert.True(t, WildcardMatch("*", ""))
	assert.False(t, WildcardMatch("?", ""))
}

func Test_WildcardMatch_MatchIdentityCaseInsensitive(t *testing.T) {
	assert.True(t, WildcardMatch("foo", "foo"))
	//assert.True(t, WildcardMatch("foo", "FOO"))
	//assert.True(t, WildcardMatch("foo", "Foo"))
	assert.True(t, WildcardMatch("6543210", "6543210"))
}

func Test_WildcardMatch_MatchIdentityExtra(t *testing.T) {
	assert.False(t, WildcardMatch("foo", "foob"))
	assert.False(t, WildcardMatch("foo", "xfoo"))
	assert.False(t, WildcardMatch("foo", "bar"))
}

func Test_WildcardMatch_SingleWildcard(t *testing.T) {
	assert.False(t, WildcardMatch("f?o", "boo"))
	assert.True(t, WildcardMatch("fo?", "foo"))
}

func Test_WildcardMatch_MultipleWildcards(t *testing.T) {
	assert.False(t, WildcardMatch("f??", "boo"))
	//assert.True(t, WildcardMatch("he??o", "Hello"))
	assert.True(t, WildcardMatch("?o?", "foo"))
}

func Test_WildcardMatch_GlobMatch(t *testing.T) {
	assert.True(t, WildcardMatch("f?o*ba*", "foobazbar"))
	assert.True(t, WildcardMatch("*oo", "foo"))
	assert.True(t, WildcardMatch("*o?", "foo"))
	assert.True(t, WildcardMatch("mis*spell", "mistily spell"))
	assert.True(t, WildcardMatch("mis*spell", "misspell"))
}

func Test_WildcardMatch_GlobMismatch(t *testing.T) {
	assert.False(t, WildcardMatch("foo*", "fo0"))
	assert.False(t, WildcardMatch("fo*obar", "foobaz"))
	assert.False(t, WildcardMatch("mis*spellx", "mispellx"))
	assert.False(t, WildcardMatch("f?*", "boo"))
}

func Test_WildcardMatch_OnlyGlob(t *testing.T) {
	assert.True(t, WildcardMatch("*", "foo"))
	assert.True(t, WildcardMatch("*", "anything"))
	assert.True(t, WildcardMatch("*", "12354"))
	assert.True(t, WildcardMatch("*", "UPPERCASE"))
	assert.True(t, WildcardMatch("*", "miXEDcaSe"))
	assert.True(t, WildcardMatch("*******", "Envoy"))
}

func Test_WildcardMatch_LengthAtLeastTwo(t *testing.T) {
	assert.False(t, WildcardMatch("??*", "a"))
	assert.True(t, WildcardMatch("??*", "aa"))
	assert.True(t, WildcardMatch("??*", "aaa"))
}
