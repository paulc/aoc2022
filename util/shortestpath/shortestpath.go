package shortestpath

import (
	"container/heap"
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

type Edge[T comparable] struct {
	To   T
	Cost float64
}

type Graph[T comparable] map[T][]Edge[T]

func (g Graph[T]) String() string {
	out := []string{}
	for k, v := range g {
		out = append(out, fmt.Sprintf("%v -> %v", k, v))
	}
	slices.Sort(out)
	return strings.Join(out, "\n")
}

// Simple Dijkstra
func (g *Graph[T]) ShortestPathSimple(start, end T) float64 {
	Q := make(map[T]struct{})
	for k, _ := range *g {
		Q[k] = struct{}{}
	}
	known := make(map[T]struct{})
	dist := make(map[T]float64)
	dist[start] = 0
	known[start] = struct{}{}
	for len(Q) > 0 {
		var u T
		min := math.Inf(1)
		for k, _ := range known {
			if dist[k] < min {
				min = dist[k]
				u = k
			}
		}
		delete(Q, u)
		delete(known, u)
		if u == end {
			break
		}
		for _, v := range (*g)[u] {
			if _, found := Q[v.To]; found {
				known[v.To] = struct{}{}
				if _, found := dist[v.To]; !found {
					dist[v.To] = math.Inf(1)
				}
				if alt := dist[u] + v.Cost; alt < dist[v.To] {
					dist[v.To] = alt
				}
			}
		}
	}
	return dist[end]
}

// Priority Queue - not thread safe
type PathQ[T comparable] struct {
	path []Edge[T]
	// We keep an index of keys->index for UpdateCost
	index map[T]int
}

func NewPathQ[T comparable]() *PathQ[T] {
	return &PathQ[T]{index: make(map[T]int)}
}

func (q PathQ[T]) Len() int {
	return len(q.path)
}

func (q PathQ[T]) Less(i, j int) bool {
	return q.path[i].Cost < q.path[j].Cost
}

func (q PathQ[T]) Swap(i, j int) {
	q.index[q.path[i].To], q.index[q.path[j].To] = q.index[q.path[j].To], q.index[q.path[i].To]
	q.path[i], q.path[j] = q.path[j], q.path[i]
}

func (q *PathQ[T]) Push(x any) {
	q.path = append(q.path, x.(Edge[T]))
	q.index[x.(Edge[T]).To] = len(q.path) - 1
}

func (q *PathQ[T]) Pop() any {
	old := q.path
	n := len(old)
	x := old[n-1]
	q.path = old[0 : n-1]
	delete(q.index, x.To)
	return x
}

func (q *PathQ[T]) UpdateCost(dest T, cost float64) {
	if i, ok := q.index[dest]; ok {
		q.path[i].Cost = cost
		heap.Fix(q, i)
	}
}

// Optimised Dijstra
func (g *Graph[T]) CalculatePaths(start T) (map[T]float64, map[T]T) {
	cost := make(map[T]float64)
	cost[start] = 0

	prev := make(map[T]T)
	pathQ := NewPathQ[T]()
	heap.Init(pathQ)

	for k, _ := range *g {
		if k != start {
			cost[k] = math.Inf(1)
			// prev[k] = ""
		}
		heap.Push(pathQ, Edge[T]{k, cost[k]})
	}

	for pathQ.Len() > 0 {
		u := heap.Pop(pathQ).(Edge[T])
		for _, v := range (*g)[u.To] {
			if alt := cost[u.To] + v.Cost; alt < cost[v.To] {
				cost[v.To] = alt
				prev[v.To] = u.To
				pathQ.UpdateCost(v.To, alt)
			}
		}
	}
	return cost, prev
}

func (g *Graph[T]) Route(start, end T) (float64, []T) {
	cost, prev := g.CalculatePaths(start)
	cur := end
	route := []T{end}
	for cur != start {
		cur = prev[cur]
		route = append(route, cur)
	}
	return cost[end], route
}
