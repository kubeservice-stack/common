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
	"fmt"

	"github.com/kubeservice-stack/common/pkg/orders"
)

// Vertex type implements a vertex of a Directed Acyclic graph or DAG.
type Vertex struct {
	ID       string
	Value    interface{}
	Parents  orders.OrderSet
	Children orders.OrderSet
}

// NewVertex creates a new vertex.
func NewVertex(id string, value interface{}) *Vertex {
	v := &Vertex{
		ID:       id,
		Parents:  orders.NewOrderedSet(),
		Children: orders.NewOrderedSet(),
		Value:    value,
	}

	return v
}

// Degree return the number of parents and children of the vertex
func (v *Vertex) Degree() int {
	return v.Parents.Size() + v.Children.Size()
}

// InDegree return the number of parents of the vertex or the number of edges
// entering on it.
func (v *Vertex) InDegree() int {
	return v.Parents.Size()
}

// OutDegree return the number of children of the vertex or the number of edges
// leaving it.
func (v *Vertex) OutDegree() int {
	return v.Children.Size()
}

// String implements stringer interface and prints an string representation
// of this instance.
func (v *Vertex) String() string {
	return fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", v.ID, v.Parents.Size(), v.Children.Size(), v.Value)
}
