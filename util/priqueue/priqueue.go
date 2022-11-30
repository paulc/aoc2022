package priqueue

type Prioritised interface {
	Priority() float64
}

type PriorityQueue[T Prioritised] []T

func (q PriorityQueue[T]) Len() int {
	return len(q)
}

func (q PriorityQueue[T]) Less(i, j int) bool {
	return (q[i]).Priority() < (q[j]).Priority()
}

func (q PriorityQueue[T]) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *PriorityQueue[T]) Push(x any) {
	*q = append(*q, x.(T))
}

func (q *PriorityQueue[T]) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}
