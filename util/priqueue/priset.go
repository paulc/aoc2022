package priqueue

import "container/heap"

type PrioritySet[T comparable] struct {
	values   []T
	members  map[T]struct{}
	priority func(T) float64
}

func NewPrioritySet[T comparable](f func(T) float64) *PrioritySet[T] {
	q := &PrioritySet[T]{priority: f, members: make(map[T]struct{})}
	heap.Init(q)
	return q
}

func (q PrioritySet[T]) Len() int {
	return len(q.values)
}

func (q PrioritySet[T]) Less(i, j int) bool {
	return q.priority(q.values[i]) < q.priority(q.values[j])
}

func (q PrioritySet[T]) Swap(i, j int) {
	q.values[i], q.values[j] = q.values[j], q.values[i]
}

func (q *PrioritySet[T]) Push(x any) {
	q.values = append(q.values, x.(T))
	q.members[x.(T)] = struct{}{}
}

func (q *PrioritySet[T]) Pop() any {
	old := q.values
	n := len(old)
	x := old[n-1]
	q.values = old[0 : n-1]
	delete(q.members, x)
	return x
}

func (q *PrioritySet[T]) Contains(key T) bool {
	_, found := q.members[key]
	return found
}
