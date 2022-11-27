package util

import (
	"container/heap"
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

type Path struct {
	To   string
	Cost float64
}

type Graph map[string][]Path

func (g Graph) String() string {
	out := []string{}
	for k, v := range g {
		out = append(out, fmt.Sprintf("%s -> %v", k, v))
	}
	slices.Sort(out)
	return strings.Join(out, "\n")
}

func (g *Graph) ShortestPath(start, end string) float64 {
	Q := make(map[string]struct{})
	for k, _ := range *g {
		Q[k] = struct{}{}
	}
	known := make(map[string]struct{})
	dist := make(map[string]float64)
	dist[start] = 0
	known[start] = struct{}{}
	for len(Q) > 0 {
		// fmt.Println("Keys:", maps.Keys(Q))
		// fmt.Println("Known:", maps.Keys(known))
		// fmt.Println("Dist:", dist)
		var u string
		min := math.Inf(1)
		for k, _ := range known {
			if dist[k] < min {
				min = dist[k]
				u = k
			}
		}
		// fmt.Println("u >>", u)
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

// XXX Not thread safe
type PathQ struct {
	path  []Path
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
	q.path = append(q.path, x.(Path))
	q.index[x.(Path).To] = len(q.path) - 1
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
		heap.Push(pathQ, Path{k, cost[k]})
	}

	for pathQ.Len() > 0 {
		u := heap.Pop(pathQ).(Path)
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
