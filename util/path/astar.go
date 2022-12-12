package path

import (
	"container/heap"
	"math"

	"github.com/paulc/aoc2022/util/priqueue"
)

// A*

type scoreMap[T comparable] map[T]float64

func (m scoreMap[T]) GetDefault(s T) float64 {
	if v, found := m[s]; found {
		return v
	} else {
		return math.Inf(1)
	}
}

type score[T comparable] struct {
	key   T
	score float64
}

func (g *Graph[T]) Astar(start, end T, h func(s T) float64) (cost float64, path []T) {
	openSet := priqueue.NewPriorityKeySet(
		func(s score[T]) float64 { return s.score },
		func(s score[T]) T { return s.key })
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
	cost = gScore[end]
	path = []T{end}
	current := end
	found := false
	for {
		current, found = cameFrom[current]
		if !found {
			break
		}
		path = append(path, current)
	}
	return
}

type AstarResult[T comparable] struct {
	end  T
	cost float64
	path []T
}

func (g *Graph[T]) AstarMultiple(start T, endList []T, h func(s T) float64) (out []AstarResult[T]) {
	openSet := priqueue.NewPriorityKeySet(
		func(s score[T]) float64 { return s.score },
		func(s score[T]) T { return s.key })
	heap.Push(openSet, score[T]{start, h(start)})
	cameFrom := map[T]T{}
	gScore := scoreMap[T]{start: 0}
	fScore := scoreMap[T]{start: h(start)}
	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(score[T])
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
	for _, end := range endList {
		result := AstarResult[T]{end: end}
		result.cost = gScore[end]
		result.path = []T{end}
		current := end
		found := false
		for {
			current, found = cameFrom[current]
			if !found {
				break
			}
			result.path = append(result.path, current)
		}
		out = append(out, result)
	}
	return
}
