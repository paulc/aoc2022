package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/grid"
	"github.com/paulc/aoc2022/util/point"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
)

type blizzard struct {
	dx, dy int
}

type puzzle struct {
	blizzards  map[point.Point][]blizzard
	w, h       int
	start, end point.Point
}

func drawBlizzard(blizzards map[point.Point][]blizzard, w, h int) string {
	g, _ := grid.NewGrid[string](0, 0, w-1, h-1)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			g.Set(point.Point{x, y}, ".")
		}
	}
	for x := 2; x < w; x++ {
		g.Set(point.Point{x, 0}, "#")
		g.Set(point.Point{w - x - 1, h - 1}, "#")
	}
	for y := 0; y < h; y++ {
		g.Set(point.Point{0, y}, "#")
		g.Set(point.Point{w - 1, y}, "#")
	}
	for k, v := range blizzards {
		if len(v) == 1 {
			switch v[0] {
			case blizzard{-1, 0}:
				g.Set(k, "<")
			case blizzard{1, 0}:
				g.Set(k, ">")
			case blizzard{0, -1}:
				g.Set(k, "^")
			case blizzard{0, 1}:
				g.Set(k, "v")
			}
		} else {
			g.Set(k, fmt.Sprintf("%d", len(v)))
		}
	}
	return g.String()
}

func parseInput(r io.Reader) (out puzzle) {
	out.blizzards = make(map[point.Point][]blizzard)
	y := 0
	for _, l := range util.Must(reader.Lines(r)) {
		out.w = len(l)
		for x, v := range strings.Split(l, "") {
			switch v {
			case "<":
				out.blizzards[point.Point{x, y}] = append(out.blizzards[point.Point{x, y}], blizzard{-1, 0})
			case ">":
				out.blizzards[point.Point{x, y}] = append(out.blizzards[point.Point{x, y}], blizzard{1, 0})
			case "^":
				out.blizzards[point.Point{x, y}] = append(out.blizzards[point.Point{x, y}], blizzard{0, -1})
			case "v":
				out.blizzards[point.Point{x, y}] = append(out.blizzards[point.Point{x, y}], blizzard{0, 1})
			}
		}
		y++
	}
	out.h = y
	out.start = point.Point{1, 0}
	out.end = point.Point{out.w - 2, out.h - 1}
	return
}

func moveBlizzards(in map[point.Point][]blizzard, w, h int) (out map[point.Point][]blizzard) {
	out = make(map[point.Point][]blizzard)
	for p, v := range in {
		for _, b := range v {
			next := point.Point{p.X + b.dx, p.Y + b.dy}
			switch {
			case next.X == 0:
				next.X = w - 2
			case next.X == w-1:
				next.X = 1
			case next.Y == 0:
				next.Y = h - 2
			case next.Y == h-1:
				next.Y = 1
			}
			out[next] = append(out[next], b)
		}
	}
	return
}

func options(p point.Point, blizzards map[point.Point][]blizzard, w, h int, start, end point.Point) (out []point.Point) {
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {0, 0}} {
		next := point.Point{p.X + v.dx, p.Y + v.dy}
		if (next.X > 0 && next.X <= w-2 && next.Y > 0 && next.Y <= h-2) || next == start || next == end {
			if _, found := blizzards[next]; !found {
				out = append(out, next)
			}
		}
	}
	return
}

func part1(input puzzle) (result int) {
	current := set.NewSetFrom[point.Point]([]point.Point{{1, 0}})
	for {
		result++
		input.blizzards = moveBlizzards(input.blizzards, input.w, input.h)
		next := set.NewSet[point.Point]()
		for p := range current {
			opts := options(p, input.blizzards, input.w, input.h, input.start, input.end)
			for _, v := range opts {
				next.Add(v)
			}
		}
		if next.Has(input.end) {
			break
		}
		current = next
	}
	return result
}

func part2(input puzzle) (result int) {
	current := set.NewSetFrom[point.Point]([]point.Point{{1, 0}})
	target := input.end
	trips := 0
	for {
		result++
		input.blizzards = moveBlizzards(input.blizzards, input.w, input.h)
		next := set.NewSet[point.Point]()
		for p := range current {
			opts := options(p, input.blizzards, input.w, input.h, input.start, input.end)
			for _, v := range opts {
				next.Add(v)
			}
		}
		if next.Has(target) {
			trips++
			if trips == 3 {
				break
			}
			if target == input.end {
				target = input.start
				current = set.NewSetFrom[point.Point]([]point.Point{input.end})
			} else {
				target = input.end
				current = set.NewSetFrom[point.Point]([]point.Point{input.start})
			}
		} else {
			current = next
		}
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
