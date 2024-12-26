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

package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVertex(t *testing.T) {
	assert := assert.New(t)
	v := NewVertex("1", nil)

	assert.NotEmpty(v.ID)
	assert.Nil(v.Value)
}

func TestVertexParents(t *testing.T) {
	assert := assert.New(t)
	v := NewVertex("1", nil)

	assert.NotNil(v.Parents)
	assert.Equal(0, v.Parents.Size())
}

func TestVertexChildren(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex("1", nil)

	assert.NotNil(v.Children)
	assert.Equal(0, v.Children.Size())
}

func TestVertexDegree(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	degree := vertex1.Degree()
	assert.Equal(0, degree)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Nil(dag1.AddEdge(vertex2, vertex3))

	degree = vertex1.Degree()
	assert.Equal(1, degree)

	degree = vertex2.Degree()
	assert.Equal(2, degree)
}

func TestVertexInDegree(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	inDegree := vertex1.InDegree()
	assert.Equal(0, inDegree)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Nil(dag1.AddEdge(vertex2, vertex3))

	inDegree = vertex1.InDegree()
	assert.Equal(0, inDegree)

	inDegree = vertex2.InDegree()
	assert.Equal(1, inDegree)
}

func TestVertexOutDegree(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	outDegree := vertex1.OutDegree()
	assert.Equal(0, outDegree)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Nil(dag1.AddEdge(vertex2, vertex3))

	outDegree = vertex1.OutDegree()
	assert.Equal(1, outDegree)

	outDegree = vertex2.OutDegree()
	assert.Equal(1, outDegree)

	outDegree = vertex3.OutDegree()
	assert.Equal(0, outDegree)
}

func TestVertexString(t *testing.T) {
	assert := assert.New(t)
	v := NewVertex("1", nil)
	vstr := v.String()

	expected := "ID: 1 - Parents: 0 - Children: 0 - Value: <nil>\n"
	assert.Equal(expected, vstr)
}

func TestVertexStringWithStringValue(t *testing.T) {
	assert := assert.New(t)
	v := NewVertex("1", "one")
	vstr := v.String()

	expected := "ID: 1 - Parents: 0 - Children: 0 - Value: one\n"
	assert.Equal(expected, vstr)
}

func TestVertexWithStringValue(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex("1", "one")

	assert.NotEmpty(v.ID)
	assert.Equal("one", v.Value)
}
