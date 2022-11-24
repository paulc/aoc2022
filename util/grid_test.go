package util

import (
	"fmt"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGridCheckBounds(t *testing.T) {

	g := NewGrid[int](0, 0, 10, 10)
	for _, v := range []Point{{0, 0}, {5, 5}, {10, 10}} {
		if !g.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
	for _, v := range []Point{{-1, 0}, {5, -5}, {10, 11}} {
		if g.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}

	g2 := NewGrid[int](-5, -5, 5, 5)
	for _, v := range []Point{{0, 0}, {5, -5}, {-5, -5}} {
		if !g2.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
	for _, v := range []Point{{-10, 0}, {5, -6}, {5, 6}} {
		if g2.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
}

type _gridData int

func (d _gridData) String() string {
	return fmt.Sprintf("%3d", d)
}

func TestGridPrint(t *testing.T) {
	g := NewGrid[_gridData](0, 0, 5, 5)
	g.Set(Point{2, 2}, 99)
	if g.Get(Point{2, 2}) != 99 {
		t.Error(g)
	}
	fmt.Println(g)
}

func TestGridAdjacent(t *testing.T) {
	g := NewGrid[_gridData](0, 0, 5, 5)

	adj := g.Adjacent(Point{2, 2})
	if !slices.Equal(adj, []Point{{1, 2}, {2, 1}, {3, 2}, {2, 3}}) {
		t.Error(adj)
	}

	adj = g.Adjacent(Point{0, 0})
	if !slices.Equal(adj, []Point{{1, 0}, {0, 1}}) {
		t.Error(adj)
	}
}
