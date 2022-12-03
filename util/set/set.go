package set

import (
	"fmt"
	"strings"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() *Set[T] {
	out := make(Set[T])
	return &out
}

func NewSetFrom[T comparable](s []T) *Set[T] {
	out := make(Set[T])
	for _, v := range s {
		out[v] = struct{}{}
	}
	return &out
}

func (s *Set[T]) String() string {
	out := []string{}
	s.Apply(func(v T) { out = append(out, fmt.Sprintf("%v", v)) })
	return fmt.Sprintf("{%s}", strings.Join(out, " "))
}

func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(*s, v)
}

func (s *Set[T]) Has(v T) bool {
	_, ok := (*s)[v]
	return ok
}

func (s1 *Set[T]) Equals(s2 *Set[T]) bool {
	if s1.Len() != s2.Len() {
		return false
	}
	for e := range *s1 {
		if !s2.Has(e) {
			return false
		}
	}
	return true
}

func (s1 *Set[T]) Intersection(s2 *Set[T]) *Set[T] {
	out := NewSet[T]()
	s1.Apply(func(v T) {
		if s2.Has(v) {
			out.Add(v)
		}
	})
	return out
}

func (s1 *Set[T]) Union(s2 *Set[T]) *Set[T] {
	out := NewSet[T]()
	s1.Apply(func(v T) { out.Add(v) })
	s2.Apply(func(v T) { out.Add(v) })
	return out
}

func (s *Set[T]) Apply(f func(T)) {
	for k, _ := range *s {
		f(k)
	}
}

func (s *Set[T]) Copy() *Set[T] {
	out := NewSet[T]()
	s.Apply(func(v T) { out.Add(v) })
	return out
}

func (s *Set[T]) Len() int {
	return len(*s)
}

func (s *Set[T]) Pop() (v T, ok bool) {
	for k, _ := range *s {
		delete(*s, k)
		return k, true
	}
	return
}

func (s *Set[T]) Iter() <-chan T {
	ch := make(chan T)
	go func() {
		for k, _ := range *s {
			ch <- k
		}
		close(ch)
	}()
	return ch
}
