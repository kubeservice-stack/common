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

func TestByteSet(t *testing.T) {
	assert := assert.New(t)

	s := Byte{}
	s2 := Byte{}
	assert.Equal(s.Len(), 0)

	s.Insert('a', 'b')
	assert.Equal(s.Len(), 2)

	s.Insert('c')
	assert.False(s.Has('d'))

	assert.True(s.Has('a'))
	s.Delete('a')
	assert.False(s.Has('a'))

	assert.True(s.HasAll('b', 'c'))
	assert.True(s.HasAny('b', 'd'))

	s2.Insert('b', 'c')
	assert.True(s.IsSuperset(s2))
	assert.True(s.Equal(s2))

	s2.Delete('d')
	assert.True(s.IsSuperset(s2))

	s2.Delete('c')
	assert.True(s.IsSuperset(s2))

	assert.Equal(s.Intersection(s2), s2)
	assert.Equal(s2.Clone(), s2)
	assert.Equal(s.Intersection(s2).List(), []byte{'b'}) // string not order
	assert.Equal(s.Intersection(s2).UnsortedList(), []byte{'b'})
	ret, ok := s.Intersection(s2).PopAny()
	assert.Equal(ret, byte('b'))
	assert.True(ok)

	s2.Insert('a', 'b', 'c')
	assert.False(s.IsSuperset(s2))
	assert.Equal(s.Union(s2), s2)
	assert.Equal(s.SymmetricDifference(s2).Len(), 1)
	assert.Equal(s.Difference(s2).Len(), 0)
	assert.Equal(s2.Difference(s).Len(), 1)

	s3 := NewByte('a', 'b', 'c')
	s4 := ByteKeySet(map[byte]interface{}{'a': 1, 'b': 1, 'c': 1})
	assert.True(s3.Equal(s4))
}
