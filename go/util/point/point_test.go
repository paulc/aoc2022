package point

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestPoint(t *testing.T) {
	p1 := Point{0, 0}
	if (p1.Move(5, 5) != Point{5, 5}) {
		t.Error(p1.Move(5, 5))
	}
	if p1.Distance(Point{5, -5}) != 10 {
		t.Error(p1.Distance(Point{5, -5}))
	}
	if adj := (Point{5, 5}).Adjacent(); !slices.Equal(adj, []Point{{4, 5}, {5, 4}, {6, 5}, {5, 6}}) {
		t.Error(adj)
	}
	if adj := (Point{5, 5}).AdjacentDiagonal(); !slices.Equal(adj, []Point{{4, 5}, {5, 4}, {6, 5}, {5, 6}, {4, 4}, {6, 4}, {4, 6}, {6, 6}}) {
		t.Error(adj)
	}
}
