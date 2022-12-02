package priqueue

import (
	"container/heap"
	"testing"

	"golang.org/x/exp/slices"
)

func TestPriorityQueue(t *testing.T) {
	q := NewPriorityQueue[int](func(i int) float64 { return float64(i) })
	heap.Push(q, 5)
	heap.Push(q, 10)
	heap.Push(q, 1)
	heap.Push(q, 3)
	if !slices.Equal(q.values, []int{1, 3, 5, 10}) {
		t.Error(*q)
	}
}
