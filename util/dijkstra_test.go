package util

import (
	"container/heap"
	"testing"

	"golang.org/x/exp/slices"
)

type _i int

func (i _i) Priority() float64 {
	return float64(i)
}

func TestPriorityQueue(t *testing.T) {
	q := &PriorityQueue[_i]{}
	heap.Init(q)
	heap.Push(q, _i(5))
	heap.Push(q, _i(10))
	heap.Push(q, _i(1))
	heap.Push(q, _i(3))
	if !slices.Equal(*q, []_i{_i(1), _i(3), _i(5), _i(10)}) {
		t.Error(*q)
	}
}
