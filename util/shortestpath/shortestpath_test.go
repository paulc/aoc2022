package shortestpath

import (
	"bytes"
	"container/heap"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/paulc/aoc2022/util/array"
	"github.com/paulc/aoc2022/util/point"
	"github.com/paulc/aoc2022/util/reader"
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

func makeGraph(r io.Reader) (*Graph[point.Point], error) {
	a, err := array.ArrayReader[int](r, array.MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	g := make(Graph[point.Point])
	a.Each(func(e array.ArrayElement[int]) {
		adj := []Edge[point.Point]{}
		for _, v := range a.Adjacent(e.X, e.Y) {
			adj = append(adj, Edge[point.Point]{point.Point{v.X, v.Y}, float64(v.Val)})
		}
		g[point.Point{e.X, e.Y}] = adj
	})
	return &g, nil
}

func makeGraphRepeat(r io.Reader) (*Graph[point.Point], error) {
	a, err := array.ArrayReader(r, array.MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	h := len(a)
	w := len(a[0])
	a2 := make(array.Array[int], h*5)
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
	g := make(Graph[point.Point])
	a2.Each(func(e array.ArrayElement[int]) {
		adj := []Edge[point.Point]{}
		for _, v := range a2.Adjacent(e.X, e.Y) {
			adj = append(adj, Edge[point.Point]{point.Point{v.X, v.Y}, float64(v.Val)})
		}
		g[point.Point{e.X, e.Y}] = adj
	})
	return &g, nil
}

func TestShortestPathSimple(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.ShortestPathSimple(point.Point{0, 0}, point.Point{9, 9})
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func TestCalculatePaths(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, prev := g.CalculatePaths(point.Point{0, 0})
	if cost[point.Point{9, 9}] != 40 {
		t.Error("cost:", cost, "\nprev:", prev)
	}
}

func TestRoute(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, _ := g.Route(point.Point{0, 0}, point.Point{9, 9})
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func makeHF(target point.Point) func(point.Point) float64 {
	return func(p point.Point) float64 {
		return float64(p.Distance(target))
	}
}

func TestAstar(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, path := g.Astar(point.Point{0, 0}, point.Point{9, 9}, makeHF(point.Point{9, 9}))
	if cost != 40 {
		t.Error("cost:", cost)
	}
	if len(path) != 19 {
		t.Error("path:", path)
	}
}

func TestGraphRepeat(t *testing.T) {
	g, err := makeGraphRepeat(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	start, end := point.Point{0, 0}, point.Point{49, 49}
	cost := g.ShortestPathSimple(start, end)
	if cost != 315 {
		t.Error("simple cost:", cost)
	}
	cost, _ = g.Route(start, end)
	if cost != 315 {
		t.Error("route cost:", cost)
	}

	cost, _ = g.Astar(start, end, makeHF(end))
	if cost != 315 {
		t.Error("astar cost:", cost)
	}
}

func BenchmarkRoute(b *testing.B) {
	r, err := reader.UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Route(point.Point{0, 0}, point.Point{99, 99})
		if cost != 602 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkAstar(b *testing.B) {
	r, err := reader.UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Astar(point.Point{0, 0}, point.Point{99, 99}, makeHF(point.Point{99, 99}))
		if cost != 602 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkRouteRepeat(b *testing.B) {
	r, err := reader.UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Route(point.Point{0, 0}, point.Point{499, 499})
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

func BenchmarkAstarRepeat(b *testing.B) {
	r, err := reader.UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost, _ := g.Astar(point.Point{0, 0}, point.Point{499, 499}, makeHF(point.Point{499, 499}))
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

/*

// No point benchmarking - 20 x slower than Astar

func BenchmarkShortestPathSimpleRepeat(b *testing.B) {
	r, err := reader.UrlOpen("testdata/input.txt")
	if err != nil {
		b.Fatal(err)
	}
	g, err := makeGraphRepeat(r)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cost := g.ShortestPathSimple(point.Point{0, 0}, point.Point{499, 499})
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

*/
