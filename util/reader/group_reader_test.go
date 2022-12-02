package reader

import (
	"bytes"
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGroupReader(t *testing.T) {

	b := bytes.NewBufferString(`1
2

3

4
5

6
7
8
9`)

	a, err := GroupReader(b, func(s string) bool { return s == "" }, func(s string) (int, error) { return strconv.Atoi(s) })
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{1, 2}, {3}, {4, 5}, {6, 7, 8, 9}}
	for i, _ := range a {
		if !slices.Equal(a[i], expected[i]) {
			t.Error(a, expected)
		}
	}
}
