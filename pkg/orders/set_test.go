/*
Copyright 2024 The KubeService-Stack Authors.

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

package orders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	assert := assert.New(t)
	s := NewOrderedSet()

	assert.True(s.Empty(), "New set expected to be empty but it is not")
	assert.Equal(0, s.Size())
}

func TestSet_Add(t *testing.T) {
	assert := assert.New(t)
	s := NewOrderedSet()
	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	s.Add("b") //overwrite
	v := FakeType{"aa"}
	s.Add(v)
	s.Add(&v)
	s.Add(true)

	vs := s.Values()
	expectedOutput := []interface{}{"e", "f", "g", "c", "d", "x", "b", "a", v, &v, true}
	assert.Equal(expectedOutput, vs)
}

func TestSet_Remove(t *testing.T) {
	assert := assert.New(t)

	s := NewOrderedSet()
	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	s.Add("b") //overwrite
	v := FakeType{"svalue"}
	s.Add(v)
	s.Add(&v)
	s.Add(true)

	s.Remove("f", "g", &v, true)

	vs := s.Values()
	expectedOutput := []interface{}{"e", "c", "d", "x", "b", "a", v}
	assert.Equal(expectedOutput, vs)

	// already removed, doesn't fails.
	s.Remove("f", "g", &v, true)
	vs = s.Values()
	expectedOutput = []interface{}{"e", "c", "d", "x", "b", "a", v}
	assert.Equal(expectedOutput, vs)
}

func TestSet_Contains(t *testing.T) {
	assert := assert.New(t)

	s := NewOrderedSet()
	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	s.Add("b") //overwrite
	v := FakeType{"svalue"}
	s.Add(v)
	s.Add(&v)
	s.Add(true)

	table := []struct {
		input          []interface{}
		expectedOutput bool
	}{
		{[]interface{}{"c", "d", &v}, true},
		{[]interface{}{"c", "d", nil}, false},
		{[]interface{}{true}, true},
		{[]interface{}{v}, true},
	}

	for _, test := range table {
		v := s.Contains(test.input...)
		assert.Equal(test.expectedOutput, v)
	}
}

func TestSet_Empty(t *testing.T) {
	assert := assert.New(t)
	s := NewOrderedSet()
	assert.True(s.Empty())

	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	assert.False(s.Empty())

	s.Remove("e", "f", "g", "c", "d", "x", "b", "a")
	assert.True(s.Empty())
}

func TestSet_Values(t *testing.T) {
	assert := assert.New(t)

	s := NewOrderedSet()
	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	s.Add("b") //overwrite
	v := FakeType{"aa"}
	s.Add(v)
	s.Add(&v)
	s.Add(true)

	vs := s.Values()
	expectedOutput := []interface{}{"e", "f", "g", "c", "d", "x", "b", "a", v, &v, true}
	assert.Equal(expectedOutput, vs)
}

func TestSet_Size(t *testing.T) {
	assert := assert.New(t)
	s := NewOrderedSet()

	assert.Equal(0, s.Size())

	s.Add("e", "f", "g", "c", "d", "x", "b", "a")
	s.Add("e", "f", "g", "c", "d", "x", "b", "a", "z") // overwrite
	assert.Equal(9, s.Size())

	s.Remove("e", "f", "g", "c", "d", "x", "b", "a", "z")
	assert.Equal(0, s.Size())
}

func TestSet_String(t *testing.T) {
	assert := assert.New(t)
	s := NewOrderedSet()

	s.Add("foo", "bar")
	expected := "[foo bar]"
	result := s.String()
	assert.Equal(expected, result)
}
