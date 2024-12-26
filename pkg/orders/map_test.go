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

type FakeType struct {
	bar string
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	assert.True(m.Empty(), "New map expected to be empty but it is not")
	assert.Equal(0, m.Size())
}

func TestMapPut(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	m.Put(1, "a") //overwrite
	m.Put(2, "b")

	structKey := FakeType{"aa"}
	structValue := FakeType{"bb"}
	m.Put(structKey, structValue)
	m.Put(&structKey, &structValue)

	m.Put(true, false)

	assert.Equal(10, m.Size())
	assert.False(m.Empty())
}

func TestMapPutOverwrite(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(1, "1")
	m.Put(1, "a") //overwrite

	v, ok := m.Get(1)
	assert.Equal("a", v)
	assert.True(ok)
}

func TestMapGet(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	m.Put(1, "a") //overwrite
	m.Put(2, "b")

	structKey := FakeType{"aa"}
	structValue := FakeType{"bb"}
	m.Put(structKey, structValue)
	m.Put(&structKey, &structValue)

	m.Put(true, false)
	assert.Equal(10, m.Size())
	assert.False(m.Empty())

	table := []struct {
		key           interface{}
		expectedValue interface{}
		expectedFound bool
	}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
		{structKey, structValue, true},
		{&structKey, &structValue, true},
		{true, false, true},
	}

	for _, test := range table {
		v, ok := m.Get(test.key)
		assert.Equal(test.expectedValue, v)
		assert.Equal(test.expectedFound, ok)
	}
}

func TestMapRemove(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()
	m.Put("bar", "foo")
	m.Put("foo", "bar")

	v, ok := m.Get("foo")
	assert.Equal("bar", v)
	assert.True(ok)

	m.Remove("foo")

	v, ok = m.Get("foo")
	assert.Nil(v)
	assert.False(ok)

	m.Remove("foo") // already removed
	v, ok = m.Get("foo")
	assert.Nil(v)
	assert.False(ok)
}

func TestMapEmpty(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()
	assert.True(m.Empty())

	m.Put("foo", "bar")
	assert.False(m.Empty())
}

func TestMap_Size(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	m.Put(1, "a") //overwrite
	m.Put(2, "b")

	structKey := FakeType{"a"}
	structValue := FakeType{"c"}
	m.Put(structKey, structValue)
	m.Put(&structKey, &structValue)

	m.Put(true, false)

	assert.Equal(10, m.Size())
	assert.False(m.Empty())
}

func TestMapKeys(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	m.Put(1, "a") //overwrite
	m.Put(2, "b")

	structKey := FakeType{"dd"}
	structValue := FakeType{"cc"}
	m.Put(structKey, structValue)
	m.Put(&structKey, &structValue)

	m.Put(true, false)

	assert.Equal(10, m.Size())
	assert.False(m.Empty())

	keys := m.Keys()
	expectedKeys := []interface{}{5, 6, 7, 3, 4, 1, 2, structKey, &structKey, true}
	assert.Equal(expectedKeys, keys)
}

func TestMapValues(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()

	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")

	m.Put(1, "a") //overwrite
	m.Put(2, "b")

	structKey := FakeType{"ad"}
	structValue := FakeType{"fsdf"}
	m.Put(structKey, structValue)
	m.Put(&structKey, &structValue)

	m.Put(true, false)

	assert.Equal(10, m.Size())
	assert.False(m.Empty())

	v := m.Values()
	expectedValues := []interface{}{"e", "f", "g", "c", "d", "a", "b", structValue, &structValue, false}
	assert.Equal(expectedValues, v)
}

func TestMapString(t *testing.T) {
	assert := assert.New(t)
	m := NewOrderedMap()
	m.Put(1, "foo")
	m.Put(2, "bar")

	expected := "[1:foo 2:bar]"
	result := m.String()
	assert.Equal(expected, result)
}
