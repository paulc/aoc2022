package util

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func intCompare(a, b int) int {
	switch {
	case a == b:
		return 0
	case a < b:
		return -1
	default:
		return 1
	}
}

func TestArrayTranspose(t *testing.T) {
	for _, v := range []struct {
		in  Array[int]
		out Array[int]
	}{
		{Array[int]{}, Array[int]{}},
		{Array[int]{{1, 2, 3}}, Array[int]{{1}, {2}, {3}}},
		{Array[int]{{1, 2, 3}, {4, 5, 6}}, Array[int]{{1, 4}, {2, 5}, {3, 6}}},
		{Array[int]{{1, 2}, {3, 4}}, Array[int]{{1, 3}, {2, 4}}},
	} {
		out := v.in.Transpose()
		if v.out.CompareFunc(out, intCompare) != 0 {
			t.Error(v, out)
		}
		if !v.in.EqualFunc(out.Transpose(), func(a, b int) bool { return a == b }) {
			t.Error(v.in, out.Transpose())
		}
	}
}

func TestArrayEach(t *testing.T) {
	a := Array[int]{{1, 2, 3}, {4, 5, 6}}
	out := []string{}
	a.Each(func(e ArrayElement[int]) { out = append(out, fmt.Sprintf("%v", e)) })
	if !slices.Equal(out, []string{"0,0: 1", "1,0: 2", "2,0: 3", "0,1: 4", "1,1: 5", "2,1: 6"}) {
		t.Error(out)
	}
}

func TestArrayAdjacent(t *testing.T) {
	a := Array[int]{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	sums := []int{}
	a.Each(func(e ArrayElement[int]) {
		sum := 0
		for _, v := range a.Adjacent(e.x, e.y) {
			sum += v.val
		}
		sums = append(sums, sum)
	})
	if !slices.Equal(sums, []int{6, 9, 8, 13, 20, 17, 12, 21, 14}) {
		t.Error(sums)
	}
}
