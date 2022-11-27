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

	heap.Push(pathQ, Path{"A", 5})
	heap.Push(pathQ, Path{"B", 3})
	heap.Push(pathQ, Path{"C", 99})
	heap.Push(pathQ, Path{"D", 1})

	out := []string{}
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Path).To)
	}

	if !slices.Equal(out, []string{"D", "B", "A", "C"}) {
		t.Error(out)
	}

	heap.Push(pathQ, Path{"A", 5})
	heap.Push(pathQ, Path{"B", 3})
	heap.Push(pathQ, Path{"C", 99})
	heap.Push(pathQ, Path{"D", 1})
	pathQ.UpdateCost("C", 2)

	out = out[:0]
	for pathQ.Len() > 0 {
		out = append(out, heap.Pop(pathQ).(Path).To)
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
	a, err := ArrayReader(r, MakeStringSplitter(""), strconv.Atoi)
	if err != nil {
		return nil, err
	}
	g := make(Graph)
	a.Each(func(e ArrayElement[int]) {
		adj := []Path{}
		for _, v := range a.Adjacent(e.x, e.y) {
			adj = append(adj, Path{fmt.Sprintf("%d:%d", v.x, v.y), float64(v.val)})
		}
		g[fmt.Sprintf("%d:%d", e.x, e.y)] = adj
	})
	return &g, nil
}

func TestShortestPath(t *testing.T) {
	g, err := makeGraph(bytes.NewBufferString(strings.TrimSpace(path_test)))
	if err != nil {
		t.Fatal(err)
	}
	sp := g.ShortestPath("0:0", "9:9")
	if sp != 40 {
		t.Error("SP:", sp)
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
	expected2 := []string{"9:9", "9:8", "8:8", "8:7", "8:6", "8:5", "7:5", "8:4", "7:3", "6:3", "6:2", "5:2", "4:2", "3:2", "2:2", "1:2", "0:2", "0:1", "0:0"}
	if !(slices.Equal(route, expected1) || slices.Equal(route, expected2)) {
		t.Error("route:   ", route)
	}
}

func TestShortestPathBig(t *testing.T) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		t.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		t.Fatal(err)
	}
	sp := g.ShortestPath("0:0", "99:99")
	if sp != 602 {
		t.Error("cost:", sp)
	}
}

func TestRouteBig(t *testing.T) {
	r, err := UrlOpen("testdata/input.txt")
	if err != nil {
		t.Fatal(err)
	}
	g, err := makeGraph(r)
	if err != nil {
		t.Fatal(err)
	}
	cost, _ := g.Route("0:0", "99:99")
	if cost != 602 {
		t.Error("cost:", cost)
	}
}
