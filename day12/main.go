package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
	"github.com/paulc/aoc2022/util/point"
	"github.com/paulc/aoc2022/util/shortestpath"
)

type Hill struct {
	start, end point.Point
	lowest     []point.Point
	graph      shortestpath.Graph[point.Point]
}

func parseInput(r io.Reader) (hill Hill) {
	a := util.Must(array.ArrayReader[byte](r, array.MakeStringSplitter(""), func(s string) (byte, error) { return s[0], nil }))
	hill.graph = make(shortestpath.Graph[point.Point])
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
		adj := []shortestpath.Edge[point.Point]{}
		for _, v := range a.Adjacent(e.X, e.Y) {
			if a[v.Y][v.X] <= (e.Val + 1) {
				adj = append(adj, shortestpath.Edge[point.Point]{point.Point{v.X, v.Y}, 1})
			}
		}
		hill.graph[point.Point{e.X, e.Y}] = adj
	})
	return
}

func part1(input Hill) (result int) {
	cost, _ := input.graph.Astar(input.start, input.end, func(p point.Point) float64 { return float64(p.Distance(input.end)) })
	return int(cost)
}

func part2(input Hill) (result int) {
	lowest := math.Inf(1)
	for _, v := range input.lowest {
		cost, _ := input.graph.Astar(v, input.end, func(p point.Point) float64 { return float64(p.Distance(input.end)) })
		if cost > 0 && cost < lowest {
			lowest = cost
		}
	}
	return int(lowest)
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
