package util

import "testing"

func TestArrayTranspose(t *testing.T) {
	for _, v := range []struct {
		in  [][]int
		out [][]int
	}{
		{[][]int{}, [][]int{}},
		{[][]int{{1, 2, 3}}, [][]int{{1}, {2}, {3}}},
		{[][]int{{1, 2, 3}, {4, 5, 6}}, [][]int{{1, 4}, {2, 5}, {3, 6}}},
	} {
		out := ArrayTranspose(v.in)
		if !ArrayEquals(out, v.out) {
			t.Error(v.in, out)
		}
		if !ArrayEquals(ArrayTranspose(out), v.in) {
			t.Error(out, ArrayTranspose(out))
		}
	}
}
