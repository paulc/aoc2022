package graph

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {

	g := make(Graph[int, string])
	v1 := NewVertex(1, "v1")
	v2 := NewVertex(2, "v2")
	v3 := NewVertex(3, "v3")
	v1.AddEdge(v2, 1)
	v2.AddEdge(v3, 2)
	g.AddVertex(v1)
	g.AddVertex(v2)
	g.AddVertex(v3)
	fmt.Println(g)
}
