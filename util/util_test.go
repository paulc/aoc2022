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

func TestGroup(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	out := [][]int{}
	for _, g := range Group(s, 3) {
		out = append(out, g)
	}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7}}
	if !slices.EqualFunc(out, expected, slices.Equal[int]) {
		t.Error(out, expected)
	}
}

func TestSplit(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	head, tail := Split(s, func(i int) bool { return i == 4 })
	if !(slices.Equal(head, []int{1, 2, 3}) && slices.Equal(tail, []int{5, 6, 7})) {
		t.Error(head, tail)
	}
}

func TestSlurpInt(t *testing.T) {
	if !slices.Equal(Must(SlurpInt(`aa1234:"|777--888kkk`)), []int{1234, 777, 888}) {
		t.Error(SlurpInt(`aa1234:"|777--888kkk`))
	}
}

func TestMap(t *testing.T) {
	m := Map([]string{"1", "2", "3"}, func(s string) int { return Must(strconv.Atoi(s)) })
	if !slices.Equal(m, []int{1, 2, 3}) {
		t.Error(m)
	}
}

func TestApply(t *testing.T) {
	count := 0
	Apply([]int{1, 2, 3, 4, 5, 6}, func(i int) { count += i })
	if count != 21 {
		t.Error(count)
	}
}

func TestFilter(t *testing.T) {
	f := Filter([]int{1, 2, 3, 4, 5, 6}, func(i int) bool { return i < 4 })
	if !slices.Equal(f, []int{1, 2, 3}) {
		t.Error(f)
	}
}

func TestReduce(t *testing.T) {
	r := Reduce([]int{1, 2, 3, 4, 5, 6}, func(a, b int) int { return a + b }, 0)
	if r != 21 {
		t.Error(r)
	}
}

func TestMin(t *testing.T) {
	if Min(2, 1) != 1 || Min(1.5, 1.4) != 1.4 {
		t.Error("Min")
	}
}

func TestMax(t *testing.T) {
	if Max(2, 1) != 2 || Max(1.5, 1.4) != 1.5 {
		t.Error("Max")
	}
}
