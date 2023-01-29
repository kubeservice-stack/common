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

func TestInt32Set(t *testing.T) {
	assert := assert.New(t)

	s := Int32{}
	s2 := Int32{}
	assert.Equal(s.Len(), 0)

	s.Insert(1, 2)
	assert.Equal(s.Len(), 2)

	s.Insert(3)
	assert.False(s.Has(4))

	assert.True(s.Has(1))
	s.Delete(1)
	assert.False(s.Has(1))

	assert.True(s.HasAll(2, 3))
	assert.True(s.HasAny(2, 4))

	s2.Insert(2, 3)
	assert.True(s.IsSuperset(s2))
	assert.True(s.Equal(s2))

	s2.Delete(4)
	assert.True(s.IsSuperset(s2))

	s2.Delete(3)
	assert.True(s.IsSuperset(s2))

	assert.Equal(s.Intersection(s2), s2)
	assert.Equal(s2.Clone(), s2)
	assert.Equal(s.Intersection(s2).List(), []int32{2}) // string not order
	assert.Equal(s.Intersection(s2).UnsortedList(), []int32{2})
	ret, ok := s.Intersection(s2).PopAny()
	assert.Equal(ret, int32(2))
	assert.True(ok)

	s2.Insert(1, 2, 3)
	assert.False(s.IsSuperset(s2))
	assert.Equal(s.Union(s2), s2)
	assert.Equal(s.SymmetricDifference(s2).Len(), 1)
	assert.Equal(s.Difference(s2).Len(), 0)
	assert.Equal(s2.Difference(s).Len(), 1)

	s3 := NewInt32(1, 2, 3)
	s4 := Int32KeySet(map[int32]interface{}{1: "", 2: "", 3: ""})
	assert.True(s3.Equal(s4))
}
