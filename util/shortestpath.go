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

type scoreQ []score

func (q scoreQ) Len() int           { return len(q) }
func (q scoreQ) Less(i, j int) bool { return q[i].score < q[j].score }
func (q scoreQ) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *scoreQ) Push(x any)        { *q = append(*q, x.(score)) }

func (q *scoreQ) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

func (g *Graph) Astar(start, end string, h func(s string) float64) float64 {
	openSet := &scoreQ{score{start, h(start)}}
	heap.Init(openSet)
	cameFrom := map[string]string{}
	gScore := scoreMap{start: 0}
	fScore := scoreMap{start: h(start)}
	for len(*openSet) > 0 {
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
				// XXX
				contains := false
				for _, v := range *openSet {
					if v.key == neighbour {
						contains = true
						break
					}
				}
				if !contains {
					heap.Push(openSet, score{neighbour, fScore[neighbour]})
				}
			}
		}
	}
	return gScore[end]
}

/*

function A_Star(start, goal, h)
    // The set of discovered nodes that may need to be (re-)expanded.
    // Initially, only the start node is known.
    // This is usually implemented as a min-heap or priority queue rather than a hash-set.
    openSet := {start}

    // For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from start
    // to n currently known.
    cameFrom := an empty map

    // For node n, gScore[n] is the cost of the cheapest path from start to n currently known.
    gScore := map with default value of Infinity
    gScore[start] := 0

    // For node n, fScore[n] := gScore[n] + h(n). fScore[n] represents our current best guess as to
    // how cheap a path could be from start to finish if it goes through n.
    fScore := map with default value of Infinity
    fScore[start] := h(start)

    while openSet is not empty
        // This operation can occur in O(Log(N)) time if openSet is a min-heap or a priority queue
        current := the node in openSet having the lowest fScore[] value
        if current = goal
            return reconstruct_path(cameFrom, current)

        openSet.Remove(current)
        for each neighbor of current
            // d(current,neighbor) is the weight of the edge from current to neighbor
            // tentative_gScore is the distance from start to the neighbor through current
            tentative_gScore := gScore[current] + d(current, neighbor)
            if tentative_gScore < gScore[neighbor]
                // This path to neighbor is better than any previous one. Record it!
                cameFrom[neighbor] := current
                gScore[neighbor] := tentative_gScore
                fScore[neighbor] := tentative_gScore + h(neighbor)
                if neighbor not in openSet
                    openSet.add(neighbor)

    // Open set is empty but goal was never reached
    return failure
*/
