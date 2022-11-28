package util

import (
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestPathQ(t *testing.T) {
	pathQ := NewPathQ()
	heap.Init(pathQ)

	heap.Push(pathQ, Edge{"A", 5})
	heap.Push(pathQ, Edge{"B", 3})
	heap.Push(pathQ, Edge{"C", 99})
	heap.Push(pathQ, Edge{"D", 1})

	out := []string{}
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Edge).To)
	}

	if !slices.Equal(out, []string{"D", "B", "A", "C"}) {
		t.Error(out)
	}

	heap.Push(pathQ, Edge{"A", 5})
	heap.Push(pathQ, Edge{"B", 3})
	heap.Push(pathQ, Edge{"C", 99})
	heap.Push(pathQ, Edge{"D", 1})
	pathQ.UpdateCost("C", 2)

	out = out[:0]
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Edge).To)
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

func makeGraph(r io.Reader) (*Graph, error) {
	a, err := ArrayReader[int](r, MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	g := make(Graph)
	a.Each(func(e ArrayElement[int]) {
		adj := []Edge{}
		for _, v := range a.Adjacent(e.x, e.y) {
			adj = append(adj, Edge{fmt.Sprintf("%d:%d", v.x, v.y), float64(v.val)})
		}
		g[fmt.Sprintf("%d:%d", e.x, e.y)] = adj
	})
	return &g, nil
}

func makeGraphRepeat(r io.Reader) (*Graph, error) {
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
	g := make(Graph)
	a2.Each(func(e ArrayElement[int]) {
		adj := []Edge{}
		for _, v := range a2.Adjacent(e.x, e.y) {
			adj = append(adj, Edge{fmt.Sprintf("%d:%d", v.x, v.y), float64(v.val)})
		}
		g[fmt.Sprintf("%d:%d", e.x, e.y)] = adj
	})
	return &g, nil
}

func TestShortestPathSimple(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.ShortestPathSimple("0:0", "9:9")
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func TestCalculatePaths(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, prev := g.CalculatePaths("0:0")
	if cost["9:9"] != 40 {
		t.Error("cost:", cost, "\nprev:", prev)
	}
}

func TestRoute(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost, route := g.Route("0:0", "9:9")
	if cost != 40 {
		t.Error("cost:", cost)
	}
	expected1 := []string{"9:9", "9:8", "8:8", "8:7", "8:6", "8:5", "7:5", "7:4", "7:3", "6:3", "6:2", "5:2", "4:2", "3:2", "2:2", "1:2", "0:2", "0:1", "0:0"}
	expected2 := []string{"9:9", "9:8", "8:8", "8:7", "8:6", "8:5", "8:4", "7:4", "7:3", "6:3", "6:2", "5:2", "4:2", "3:2", "2:2", "1:2", "0:2", "0:1", "0:0"}
	if !(slices.Equal(route, expected1) || slices.Equal(route, expected2)) {
		t.Error("route:", route, cost)
		t.Error("e1   :", expected1)
		t.Error("e2   :", expected2)
	}
}

func makeHF(target string) func(string) float64 {
	var p Point
	fmt.Sscanf(target, "%d:%d", &p.X, &p.Y)
	return func(s string) float64 {
		/*
			var p1 Point
			fmt.Sscanf(s, "%d:%d", &p1.X, &p1.X)
			return float64(p.Distance(p1))
		*/
		return 1
	}
}

func TestAstar(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.Astar("0:0", "9:9", makeHF("9:9"))
	if cost != 40 {
		t.Error("cost:", cost)
	}
}

func TestGraphRepeat(t *testing.T) {
	g, err := makeGraphRepeat(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	cost := g.ShortestPathSimple("0:0", "49:49")
	if cost != 315 {
		t.Error("simple cost:", cost)
	}
	cost, _ = g.Route("0:0", "49:49")
	if cost != 315 {
		t.Error("route cost:", cost)
	}

	cost = g.Astar("0:0", "49:49", makeHF("49:49"))
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
		cost, _ := g.Route("0:0", "99:99")
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
		cost := g.Astar("0:0", "99:99", makeHF("99:99"))
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
		cost, _ := g.Route("0:0", "499:499")
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
		cost := g.Astar("0:0", "499:499", makeHF("499:499"))
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}

/*
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
		cost := g.ShortestPathSimple("0:0", "499:499")
		if cost != 2935 {
			b.Error("cost:", cost)
		}
	}
}
*/
