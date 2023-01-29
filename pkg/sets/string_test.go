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

func TestStringSet(t *testing.T) {
	assert := assert.New(t)

	s := String{}
	s2 := String{}
	assert.Equal(s.Len(), 0)

	s.Insert("aa", "bb")
	assert.Equal(s.Len(), 2)

	s.Insert("cc")
	assert.False(s.Has("dd"))

	assert.True(s.Has("aa"))
	s.Delete("aa")
	assert.False(s.Has("aa"))

	assert.True(s.HasAll("bb", "cc"))
	assert.True(s.HasAny("bb", "dd"))

	s2.Insert("bb", "cc")
	assert.True(s.IsSuperset(s2))
	assert.True(s.Equal(s2))

	s2.Delete("dd")
	assert.True(s.IsSuperset(s2))

	s2.Delete("cc")
	assert.True(s.IsSuperset(s2))

	assert.Equal(s.Intersection(s2), s2)
	assert.Equal(s2.Clone(), s2)
	assert.Equal(s.Intersection(s2).List(), []string{"bb"}) // string not order
	assert.Equal(s.Intersection(s2).UnsortedList(), []string{"bb"})
	ret, ok := s.Intersection(s2).PopAny()
	assert.Equal(ret, "bb")
	assert.True(ok)

	s2.Insert("aa", "bb", "cc")
	assert.False(s.IsSuperset(s2))
	assert.Equal(s.Union(s2), s2)
	assert.Equal(s.SymmetricDifference(s2).Len(), 1)
	assert.Equal(s.Difference(s2).Len(), 0)
	assert.Equal(s2.Difference(s).Len(), 1)

	s3 := NewString("aa", "bb", "cc")
	s4 := StringKeySet(map[string]interface{}{"aa": "aa", "bb": "bb", "cc": "cc"})
	assert.True(s3.Equal(s4))
}
