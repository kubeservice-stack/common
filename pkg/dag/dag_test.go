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

func TestDAG(t *testing.T) {
	assert := assert.New(t)
	d := NewDAG()

	assert.Equal(0, d.Order())
}

func TestDAGAddVertex(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Equal(1, dag1.Order())
}

func TestDAGDeleteVertex(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Equal(1, dag1.Order())

	assert.Nil(dag1.DeleteVertex(vertex1))
	assert.Equal(0, dag1.Order())

	// delete nil
	assert.NotNil(dag1.DeleteVertex(vertex1))
}

func TestDAGAddEdge(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", "two")

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))
}

func TestDAGAddEdgeFailsVertextDontExist(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.NotNil(dag1.AddEdge(vertex3, vertex2))
	assert.NotNil(dag1.AddEdge(vertex2, vertex3))
}

func TestDAGAddEdgeFailsAlreadyExists(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.NotNil(dag1.AddEdge(vertex1, vertex2))
}

func TestDAGDeleteEdge(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	assert.Equal(1, dag1.Size())
	assert.Nil(dag1.DeleteEdge(vertex1, vertex2))
	assert.Equal(0, dag1.Size())
}

func TestDAGGetVertex(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", "one")
	vertex2 := NewVertex("2", 2)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))

	v1, err := dag1.GetVertex("1")
	assert.Nil(err)
	assert.Equal("one", v1.Value)
	v2, err := dag1.GetVertex("2")
	assert.Nil(err)
	assert.Equal(2, v2.Value)
}

func TestDAGOrder(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	assert.Equal(0, dag1.Order())

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))

	assert.Equal(3, dag1.Order())
}

func TestDAGSize(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()
	assert.Equal(0, dag1.Size())

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)
	vertex4 := NewVertex("4", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))
	assert.Nil(dag1.AddVertex(vertex4))

	assert.Equal(0, dag1.Size())

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Nil(dag1.AddEdge(vertex2, vertex3))
	assert.NotNil(dag1.AddEdge(vertex2, vertex3))

	assert.Equal(2, dag1.Size())
}

func TestDAGSinkVertices(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	assert.Len(dag1.SinkVertices(), 0)

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))

	assert.Len(dag1.SinkVertices(), 2)

	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	assert.Len(dag1.SinkVertices(), 1)
}

func TestDAGSourceVertices(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	assert.Len(dag1.SinkVertices(), 0)

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))

	assert.Len(dag1.SinkVertices(), 2)

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Len(dag1.SinkVertices(), 1)
}

func TestDAGSuccessors(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	successors, err := dag1.Successors(vertex1)
	assert.Nil(err)
	assert.Len(successors, 1)
	assert.Equal("2", successors[0].ID)
}

func TestDAGSuccessorsVertexNotFound(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	successors, err := dag1.Successors(vertex3)
	assert.NotNil(err)
	assert.Len(successors, 0)
}

func TestDAGPredecessors(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	predecessors, err := dag1.Predecessors(vertex2)
	assert.Nil(err)
	assert.Len(predecessors, 1)
	assert.Equal("1", predecessors[0].ID)
}

func TestDAGPredecessorsVertexNotFound(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddEdge(vertex1, vertex2))

	predecessors, err := dag1.Predecessors(vertex3)
	assert.NotNil(err)
	assert.Len(predecessors, 0)
}

func TestDAGString(t *testing.T) {
	assert := assert.New(t)
	dag1 := NewDAG()

	vertex1 := NewVertex("1", nil)
	vertex2 := NewVertex("2", nil)
	vertex3 := NewVertex("3", nil)
	vertex4 := NewVertex("4", nil)

	assert.Nil(dag1.AddVertex(vertex1))
	assert.Nil(dag1.AddVertex(vertex2))
	assert.Nil(dag1.AddVertex(vertex3))
	assert.Nil(dag1.AddVertex(vertex4))

	assert.Equal(0, dag1.Size())

	assert.Nil(dag1.AddEdge(vertex1, vertex2))
	assert.Nil(dag1.AddEdge(vertex2, vertex3))

	expected := `DAG Vertices: 4 - Edges: 2
Vertices:
ID: 1 - Parents: 0 - Children: 1 - Value: <nil>
ID: 2 - Parents: 1 - Children: 1 - Value: <nil>
ID: 3 - Parents: 1 - Children: 0 - Value: <nil>
ID: 4 - Parents: 0 - Children: 0 - Value: <nil>
`
	assert.Equal(expected, dag1.String())
}
