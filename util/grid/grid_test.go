package grid

import (
	"fmt"
	"testing"

	"github.com/paulc/aoc2022/util/point"
	"golang.org/x/exp/slices"
)

type _gridData int

func (d _gridData) String() string {
	return fmt.Sprintf("%3d", d)
}

func TestGridSetGetPrint(t *testing.T) {
	for _, v := range []struct{ x0, y0, x1, y1 int }{{0, 0, 5, 5}, {-2, -2, 2, 2}} {
		g, err := NewGrid[point.Point](v.x0, v.y0, v.x1, v.y1)
		if err != nil {
			t.Fatal(err)
		}
		for x := g.X0; x <= g.X1; x++ {
			for y := g.Y0; y <= g.Y1; y++ {
				g.Set(point.Point{x, y}, point.Point{x, y})
			}
		}
		t.Log("\n", g)
		for x := g.X0; x <= g.X1; x++ {
			for y := g.Y0; y <= g.Y1; y++ {
				if !(g.Get(point.Point{x, y}) == point.Point{x, y}) {
					t.Errorf("%v != %v", point.Point{x, y}, g.Get(point.Point{x, y}))
				}
			}
		}
	}
}

func TestGridInvalid(t *testing.T) {
	for _, v := range []struct{ x0, y0, x1, y1 int }{{0, 0, 0, 0}, {0, 0, -5, 5}, {-2, -2, 2, -2}} {
		_, err := NewGrid[point.Point](v.x0, v.y0, v.x1, v.y1)
		if err == nil {
			t.Error("Expected error", v)
		}
	}
}

func TestGridCheckBounds(t *testing.T) {

	g, err := NewGrid[int](0, 0, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []point.Point{{0, 0}, {5, 5}, {10, 10}} {
		if !g.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
	for _, v := range []point.Point{{-1, 0}, {5, -5}, {10, 11}} {
		if g.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}

	g2, err := NewGrid[int](-5, -5, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []point.Point{{0, 0}, {5, -5}, {-5, -5}} {
		if !g2.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
	for _, v := range []point.Point{{-10, 0}, {5, -6}, {5, 6}} {
		if g2.CheckBounds(v) {
			t.Error("CheckBounds:", v)
		}
	}
}

func TestGridAdjacent(t *testing.T) {
	g, err := NewGrid[_gridData](0, 0, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []struct {
		p   point.Point
		adj []point.Point
	}{
		{point.Point{2, 2}, []point.Point{{1, 2}, {2, 1}, {3, 2}, {2, 3}}},
		{point.Point{1, 2}, []point.Point{{0, 2}, {1, 1}, {2, 2}, {1, 3}}},
		{point.Point{0, 0}, []point.Point{{1, 0}, {0, 1}}},
		{point.Point{5, 5}, []point.Point{{4, 5}, {5, 4}}},
	} {
		if !slices.Equal(g.Adjacent(v.p), v.adj) {
			t.Error(v, g.Adjacent(v.p))
		}
	}
}

func TestGridAdjacentWrap(t *testing.T) {
	g, err := NewGrid[_gridData](0, 0, 5, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []struct {
		p   point.Point
		adj []point.Point
	}{
		{point.Point{2, 2}, []point.Point{{1, 2}, {2, 1}, {3, 2}, {2, 3}}},
		{point.Point{0, 0}, []point.Point{{5, 0}, {0, 5}, {1, 0}, {0, 1}}},
		{point.Point{5, 5}, []point.Point{{4, 5}, {5, 4}, {0, 5}, {5, 0}}},
	} {
		if !slices.Equal(g.AdjacentWrap(v.p), v.adj) {
			t.Error(v, g.AdjacentWrap(v.p))
		}
	}
}

func TestGridMove(t *testing.T) {
	g, err := NewGrid[_gridData](0, 0, 5, 5)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range []struct {
		p      point.Point
		dx, dy int
		p2     point.Point
	}{
		{point.Point{0, 0}, -1, -1, point.Point{5, 5}},
		{point.Point{0, 0}, -2, -2, point.Point{4, 4}},
		{point.Point{2, 2}, -4, 1, point.Point{4, 3}},
		{point.Point{5, 5}, 2, -2, point.Point{1, 3}},
		{point.Point{0, 0}, 20, 20, point.Point{2, 2}},
	} {
		if g.Move(v.p, v.dx, v.dy) != v.p2 {
			t.Error(v, "::", g.Move(v.p, v.dx, v.dy))
		}
	}
}

func TestGridMove2(t *testing.T) {
	g, err := NewGrid[_gridData](-2, -2, 2, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range []struct {
		p      point.Point
		dx, dy int
		p2     point.Point
	}{
		{point.Point{0, 0}, -1, -1, point.Point{-1, -1}},
		{point.Point{0, 0}, 3, 3, point.Point{-2, -2}},
		{point.Point{0, 0}, -3, -3, point.Point{2, 2}},
		{point.Point{0, 0}, 10, 10, point.Point{0, 0}},
		{point.Point{0, 0}, -20, -20, point.Point{0, 0}},
	} {
		if g.Move(v.p, v.dx, v.dy) != v.p2 {
			t.Error(v, "::", g.Move(v.p, v.dx, v.dy))
		}
	}

}
