package util

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

type scoreMap[T comparable] map[T]float64

func (m scoreMap[T]) GetDefault(s T) float64 {
	if v, found := m[s]; found {
		return v
	} else {
		return math.Inf(1)
	}
}

// A*

type score[T comparable] struct {
	key   T
	score float64
}

type scoreQ[T comparable] struct {
	scores  []score[T]
	members map[T]struct{}
}

func newScoreQ[T comparable]() *scoreQ[T] {
	q := &scoreQ[T]{scores: []score[T]{}, members: make(map[T]struct{})}
	heap.Init(q)
	return q
}

func (q scoreQ[T]) Len() int           { return len(q.scores) }
func (q scoreQ[T]) Less(i, j int) bool { return q.scores[i].score < q.scores[j].score }
func (q scoreQ[T]) Swap(i, j int)      { q.scores[i], q.scores[j] = q.scores[j], q.scores[i] }

func (q *scoreQ[T]) Push(x any) {
	q.members[x.(score[T]).key] = struct{}{}
	q.scores = append(q.scores, x.(score[T]))
}

func (q *scoreQ[T]) Pop() any {
	old := q.scores
	n := len(old)
	x := old[n-1]
	q.scores = old[0 : n-1]
	delete(q.members, x.key)
	return x
}

func (q *scoreQ[T]) Contains(key T) bool {
	_, found := q.members[key]
	return found
}

func (g *Graph[T]) Astar(start, end T, h func(s T) float64) float64 {
	openSet := newScoreQ[T]()
	heap.Push(openSet, score[T]{start, h(start)})
	cameFrom := map[T]T{}
	gScore := scoreMap[T]{start: 0}
	fScore := scoreMap[T]{start: h(start)}
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(score[T])
		if current.key == end {
			break
		}
		for _, v := range (*g)[current.key] {
			neighbour := v.To
			if tentative_gScore := gScore.GetDefault(current.key) + v.Cost; tentative_gScore < gScore.GetDefault(neighbour) {
				cameFrom[neighbour] = current.key
				gScore[neighbour] = tentative_gScore
				fScore[neighbour] = tentative_gScore + h(neighbour)
				if !openSet.Contains(neighbour) {
					heap.Push(openSet, score[T]{neighbour, fScore[neighbour]})
				}
			}
		}
	}
	return gScore[end]
}
