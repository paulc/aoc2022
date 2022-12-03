package util

import (
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()

	if Must(strconv.Atoi("99")) != 99 {
		t.Error("99")
	}

	// Should panic
	Must(strconv.Atoi("AAA"))
}

func TestTake(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	out := [][]int{}
	for _, g := range Take(s, 3) {
		out = append(out, g)
	}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	if !slices.EqualFunc(out, expected, slices.Equal[int]) {
		t.Error(out, expected)
	}
}
