package util

import (
	"container/heap"
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

type Edge struct {
	To   string
	Cost float64
}

type Graph map[string][]Edge

func (g Graph) String() string {
	out := []string{}
	for k, v := range g {
		out = append(out, fmt.Sprintf("%s -> %v", k, v))
	}
	slices.Sort(out)
	return strings.Join(out, "\n")
}

// Simple Dijkstra
func (g *Graph) ShortestPathSimple(start, end string) float64 {
	Q := make(map[string]struct{})
	for k, _ := range *g {
		Q[k] = struct{}{}
	}
	known := make(map[string]struct{})
	dist := make(map[string]float64)
	dist[start] = 0
	known[start] = struct{}{}
	for len(Q) > 0 {
		var u string
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
type PathQ struct {
	path []Edge
	// We keep an index of keys->index for UpdateCost
	index map[string]int
}

func NewPathQ() *PathQ {
	return &PathQ{index: make(map[string]int)}
}

func (q PathQ) Len() int {
	return len(q.path)
}

func (q PathQ) Less(i, j int) bool {
	return q.path[i].Cost < q.path[j].Cost
}

func (q PathQ) Swap(i, j int) {
	q.index[q.path[i].To], q.index[q.path[j].To] = q.index[q.path[j].To], q.index[q.path[i].To]
	q.path[i], q.path[j] = q.path[j], q.path[i]
}

func (q *PathQ) Push(x any) {
	q.path = append(q.path, x.(Edge))
	q.index[x.(Edge).To] = len(q.path) - 1
}

func (q *PathQ) Pop() any {
	old := q.path
	n := len(old)
	x := old[n-1]
	q.path = old[0 : n-1]
	delete(q.index, x.To)
	return x
}

func (q *PathQ) UpdateCost(dest string, cost float64) {
	if i, ok := q.index[dest]; ok {
		q.path[i].Cost = cost
		heap.Fix(q, i)
	}
}

// Optimised Dijstra
func (g *Graph) CalculatePaths(start string) (map[string]float64, map[string]string) {
	cost := make(map[string]float64)
	cost[start] = 0

	prev := make(map[string]string)
	pathQ := NewPathQ()
	heap.Init(pathQ)

	for k, _ := range *g {
		if k != start {
			cost[k] = math.Inf(1)
			prev[k] = ""
		}
		heap.Push(pathQ, Edge{k, cost[k]})
	}

	for pathQ.Len() > 0 {
		u := heap.Pop(pathQ).(Edge)
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

func (g *Graph) Route(start, end string) (float64, []string) {
	cost, prev := g.CalculatePaths(start)
	cur := end
	route := []string{end}
	for cur != start {
		cur = prev[cur]
		route = append(route, cur)
	}
	return cost[end], route
}

type scoreMap map[string]float64

func (m scoreMap) GetDefault(s string) float64 {
	if v, found := m[s]; found {
		return v
	} else {
		return math.Inf(1)
	}
}

// A*

type score struct {
	key   string
	score float64
}

type scoreQ struct {
	scores  []score
	members map[string]struct{}
}

func newScoreQ() *scoreQ {
	q := &scoreQ{scores: []score{}, members: make(map[string]struct{})}
	heap.Init(q)
	return q
}

func (q scoreQ) Len() int           { return len(q.scores) }
func (q scoreQ) Less(i, j int) bool { return q.scores[i].score < q.scores[j].score }
func (q scoreQ) Swap(i, j int)      { q.scores[i], q.scores[j] = q.scores[j], q.scores[i] }

func (q *scoreQ) Push(x any) {
	q.members[x.(score).key] = struct{}{}
	q.scores = append(q.scores, x.(score))
}

func (q *scoreQ) Pop() any {
	old := q.scores
	n := len(old)
	x := old[n-1]
	q.scores = old[0 : n-1]
	delete(q.members, x.key)
	return x
}

func (q *scoreQ) Contains(key string) bool {
	_, found := q.members[key]
	return found
}

func (g *Graph) Astar(start, end string, h func(s string) float64) float64 {
	openSet := newScoreQ()
	heap.Push(openSet, score{start, h(start)})
	cameFrom := map[string]string{}
	gScore := scoreMap{start: 0}
	fScore := scoreMap{start: h(start)}
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(score)
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
					heap.Push(openSet, score{neighbour, fScore[neighbour]})
				}
			}
		}
	}
	return gScore[end]
}
