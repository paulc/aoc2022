package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
	"github.com/paulc/aoc2022/util/path"
	"github.com/paulc/aoc2022/util/point"
	"golang.org/x/exp/slices"
)

type Hill struct {
	start, end point.Point
	lowest     []point.Point
	graph      path.Graph[point.Point]
	inverted   path.Graph[point.Point]
}

func parseInput(r io.Reader) (hill Hill) {
	a := util.Must(array.ArrayReader[byte](r, array.MakeStringSplitter(""), func(s string) (byte, error) { return s[0], nil }))
	hill.graph = make(path.Graph[point.Point])
	hill.inverted = make(path.Graph[point.Point])
	a.Each(func(p array.ArrayElement[byte]) {
		if p.Val == 'S' {
			hill.start = point.Point{p.X, p.Y}
			a[p.Y][p.X] = 'a'
			hill.lowest = append(hill.lowest, point.Point{p.X, p.Y})
		} else if p.Val == 'E' {
			hill.end = point.Point{p.X, p.Y}
			a[p.Y][p.X] = 'z'
		} else if p.Val == 'a' {
			hill.lowest = append(hill.lowest, point.Point{p.X, p.Y})
		}
	})
	a.Each(func(e array.ArrayElement[byte]) {
		adj := []path.Edge[point.Point]{}
		adj_inv := []path.Edge[point.Point]{}
		for _, v := range a.Adjacent(e.X, e.Y) {
			if a[v.Y][v.X] <= (e.Val + 1) {
				adj = append(adj, path.Edge[point.Point]{point.Point{v.X, v.Y}, 1})
			}
			if a[v.Y][v.X] >= (e.Val - 1) {
				adj_inv = append(adj_inv, path.Edge[point.Point]{point.Point{v.X, v.Y}, 1})
			}
		}
		hill.graph[point.Point{e.X, e.Y}] = adj
		hill.inverted[point.Point{e.X, e.Y}] = adj_inv
	})
	return
}

func part1(input Hill) (result int) {
	cost, _ := input.graph.Astar(input.start, input.end, func(p point.Point) float64 { return float64(p.Distance(input.end)) })
	return int(cost)
}

func part2(input Hill) (result int) {
	results := util.Map(input.inverted.AstarMultiple(input.end, input.lowest, func(p point.Point) float64 { return float64(p.Distance(input.end)) }),
		func(r path.AstarResult[point.Point]) float64 { return r.Cost })
	slices.Sort(results)
	return int(results[0])
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
