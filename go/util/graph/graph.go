package graph

import "fmt"

type Vertex[T1 comparable, T2 any] struct {
	Key   T1
	Value T2
	Edges []Edge[T1, T2]
}

func (v *Vertex[T1, T2]) String() string {
	return fmt.Sprintf("k=%v,v=%v", v.Key, v.Value)
}

func (v *Vertex[T1, T2]) AddEdge(to *Vertex[T1, T2], cost float64) {
	v.Edges = append(v.Edges, Edge[T1, T2]{to, cost})
}

func NewVertex[T1 comparable, T2 any](key T1, value T2) *Vertex[T1, T2] {
	return &Vertex[T1, T2]{Key: key, Value: value}
}

type Edge[T1 comparable, T2 any] struct {
	To   *Vertex[T1, T2]
	Cost float64
}

type Graph[T1 comparable, T2 any] map[T1]*Vertex[T1, T2]

func (g Graph[T1, T2]) AddVertex(v *Vertex[T1, T2]) {
	g[v.Key] = v
}
