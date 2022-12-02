package priqueue

import "container/heap"

type PriorityQueue[T any] struct {
	values   []T
	priority func(T) float64
}

func NewPriorityQueue[T any](f func(T) float64) *PriorityQueue[T] {
	q := &PriorityQueue[T]{priority: f}
	heap.Init(q)
	return q
}

func (q PriorityQueue[T]) Len() int {
	return len(q.values)
}

func (q PriorityQueue[T]) Less(i, j int) bool {
	return q.priority(q.values[i]) < q.priority(q.values[j])
}

func (q PriorityQueue[T]) Swap(i, j int) {
	q.values[i], q.values[j] = q.values[j], q.values[i]
}

func (q *PriorityQueue[T]) Push(x any) {
	q.values = append(q.values, x.(T))
}

func (q *PriorityQueue[T]) Pop() any {
	old := q.values
	n := len(old)
	x := old[n-1]
	q.values = old[0 : n-1]
	return x
}
