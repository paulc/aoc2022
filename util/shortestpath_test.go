package util

import (
	"bytes"
	"container/heap"
	"io"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestPathQ(t *testing.T) {
	pathQ := NewPathQ[string]()
	heap.Init(pathQ)

	heap.Push(pathQ, Edge[string]{"A", 5})
	heap.Push(pathQ, Edge[string]{"B", 3})
	heap.Push(pathQ, Edge[string]{"C", 99})
	heap.Push(pathQ, Edge[string]{"D", 1})

	out := []string{}
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Edge[string]).To)
	}

	if !slices.Equal(out, []string{"D", "B", "A", "C"}) {
		t.Error(out)
	}

	heap.Push(pathQ, Edge[string]{"A", 5})
	heap.Push(pathQ, Edge[string]{"B", 3})
	heap.Push(pathQ, Edge[string]{"C", 99})
	heap.Push(pathQ, Edge[string]{"D", 1})
	pathQ.UpdateCost("C", 2)

	out = out[:0]
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Edge[string]).To)
	}

	if !slices.Equal(out, []string{"D", "C", "B", "A"}) {
		t.Error(out)
	}
}

// From aoc2021/day15
const path_test = `
1163751742
1381373672
2136511328
3694931569
7463417111
1319128137
1359912421
3125421639
1293138521
2311944581
`

func makeGraph(r io.Reader) (*Graph[Point], error) {
	a, err := ArrayReader[int](r, MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	g := make(Graph[Point])
	a.Each(func(e ArrayElement[int]) {
		adj := []Edge[Point]{}
		for _, v := range a.Adjacent(e.x, e.y) {
			adj = append(adj, Edge[Point]{Point{v.x, v.y}, float64(v.val)})
		}
		g[Point{e.x, e.y}] = adj
	})
	return &g, nil
}

func makeGraphRepeat(r io.Reader) (*Graph[Point], error) {
	a, err := ArrayReader(r, MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	h := len(a)
	w := len(a[0])
	a2 := make(Array[int], h*5)
	for y := 0; y < h; y++ {
		for j := 0; j < 5; j++ {
			a2[y+h*j] = make([]int, w*5)
			for x := 0; x < w; x++ {
				for i := 0; i < 5; i++ {
					v := a[y][x] + i + j
					if v > 9 {
						v -= 9
					}
					a2[y+h*j][x+w*i] = v
				}
			}
		}
	}
	g := make(Graph[Point])
	a2.Each(func(e ArrayElement[int]) {
		adj := []Edge[Point]{}
		for _, v := range a2.Adjacent(e.x, e.y) {
			adj = append(adj, Edge[Point]{Point{v.x, v.y}, float64(v.val)})
		}
		g[Point{e.x, e.y}] = adj
	})
	return &g, nil
}

func TestShortestPathSimple(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.ShortestPathSimple(Point{0, 0}, Point{9, 9})
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func TestCalculatePaths(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, prev := g.CalculatePaths(Point{0, 0})
	if cost[Point{9, 9}] != 40 {
		t.Error("cost:", cost, "\nprev:", prev)
	}
}

func TestRoute(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, _ := g.Route(Point{0, 0}, Point{9, 9})
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func makeHF(target Point) func(Point) float64 {
	return func(p Point) float64 {
		return float64(p.Distance(target))
	}
}

func TestAstar(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.Astar(Point{0, 0}, Point{9, 9}, makeHF(Point{9, 9}))
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func TestGraphRepeat(t *testing.T) {
	g, err := makeGraphRepeat(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	start, end := Point{0, 0}, Point{49, 49}
	cost := g.ShortestPathSimple(start, end)
	if cost != 315 {
		t.Error("simple cost:", cost)
	}
	cost, _ = g.Route(start, end)
	if cost != 315 {
		t.Error("route cost:", cost)
	}

	cost = g.Astar(start, end, makeHF(end))
	if cost != 315 {
		t.Error("astar cost:", cost)
	}
}

func BenchmarkRoute(b *testing.B) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Route(Point{0, 0}, Point{99, 99})
		if cost != 602 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkAstar(b *testing.B) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost := g.Astar(Point{0, 0}, Point{99, 99}, makeHF(Point{99, 99}))
		if cost != 602 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkRouteRepeat(b *testing.B) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Route(Point{0, 0}, Point{499, 499})
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkAstarRepeat(b *testing.B) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost := g.Astar(Point{0, 0}, Point{499, 499}, makeHF(Point{499, 499}))
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

/*

// No point benchmarking - 20 x slower than Astar

func BenchmarkShortestPathSimpleRepeat(b *testing.B) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost := g.ShortestPathSimple(Point{0, 0}, Point{499, 499})
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

*/
