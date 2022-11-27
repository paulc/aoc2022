package util

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Array[T any] [][]T

type ArrayElement[T any] struct {
	x, y int
	val  T
}

func (e ArrayElement[T]) String() string {
	return fmt.Sprintf("%d,%d: %v", e.x, e.y, e.val)
}

func (a Array[T]) CompareFunc(b Array[T], f func(a, b T) int) int {
	return slices.CompareFunc(a, b, func(s1, s2 []T) int { return slices.CompareFunc(s1, s2, func(e1, e2 T) int { return f(e1, e2) }) })
}

func (a Array[T]) EqualFunc(b Array[T], f func(a, b T) bool) bool {
	return slices.EqualFunc(a, b, func(s1, s2 []T) bool { return slices.EqualFunc(s1, s2, func(e1, e2 T) bool { return f(e1, e2) }) })
}

func (a Array[T]) String() string {
	h := len(a)
	if h == 0 {
		return ""
	}
	w := len(a[0])
	rows := make([]string, h)
	for y := 0; y < h; y++ {
		line := make([]string, w)
		for x := 0; x < w; x++ {
			line[x] = fmt.Sprintf("%v", a[y][x])
		}
		rows[y] = strings.Join(line, " ")
	}
	return strings.Join(rows, "\n")
}

func (a Array[T]) Transpose() Array[T] {
	if len(a) == 0 {
		return Array[T]{}
	}
	h := len(a)
	w := len(a[0])
	out := make(Array[T], w)
	for x := 0; x < w; x++ {
		out[x] = make([]T, h)
	}
	for x := 0; x < h; x++ {
		for y := 0; y < w; y++ {
			out[y][x] = a[x][y]
		}
	}
	return out
}

func (a Array[T]) Each(f func(ArrayElement[T])) {
	if len(a) == 0 {
		return
	}
	h := len(a)
	w := len(a[0])
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			f(ArrayElement[T]{x, y, a[y][x]})
		}
	}
}

func (a Array[T]) Adjacent(x, y int) (out []ArrayElement[T]) {
	if len(a) == 0 {
		return
	}
	h := len(a)
	w := len(a[0])
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
		x1 := x + v.dx
		y1 := y + v.dy
		if x1 >= 0 && x1 < w && y1 >= 0 && y1 < h {
			out = append(out, ArrayElement[T]{x1, y1, a[y1][x1]})
		}
	}
	return
}
