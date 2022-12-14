package priqueue

import (
	"container/heap"
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func TestPrioritySet(t *testing.T) {

	q := NewPrioritySet(func(i int) float64 { return float64(i) })
	heap.Push(q, 5)
	heap.Push(q, 10)
	heap.Push(q, 1)
	heap.Push(q, 3)

	if !slices.Equal(q.values, []int{1, 3, 5, 10}) {
		t.Error(*q)
	}

	for _, v := range []struct {
		v int
		b bool
	}{{5, true}, {7, false}} {
		if q.Contains(v.v) != v.b {
			t.Error(v.v, v.b)
		}
	}

	if heap.Pop(q) != 1 || heap.Pop(q) != 3 {
		t.Error("Pop")
	}

	heap.Push(q, 99)

	for _, v := range []struct {
		v int
		b bool
	}{{1, false}, {3, false}, {5, true}, {99, true}} {
		if q.Contains(v.v) != v.b {
			t.Error(v.v, v.b)
		}
	}
}

func TestPriorityKeySet(t *testing.T) {

	q := NewPriorityKeySet(func(i int) float64 { return float64(i) }, func(i int) string { return fmt.Sprintf("%03d", i) })
	heap.Push(q, 5)
	heap.Push(q, 10)
	heap.Push(q, 1)
	heap.Push(q, 3)

	if !slices.Equal(q.values, []int{1, 3, 5, 10}) {
		t.Error(*q)
	}

	for _, v := range []struct {
		v string
		b bool
	}{{"005", true}, {"007", false}} {
		if q.Contains(v.v) != v.b {
			t.Error(v.v, v.b)
		}
	}

	if heap.Pop(q) != 1 || heap.Pop(q) != 3 {
		t.Error("Pop")
	}

	heap.Push(q, 99)

	for _, v := range []struct {
		v string
		b bool
	}{{"001", false}, {"003", false}, {"005", true}, {"099", true}} {
		if q.Contains(v.v) != v.b {
			t.Error(v.v, v.b)
		}
	}
}
