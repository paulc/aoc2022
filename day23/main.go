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
	"golang.org/x/exp/slices"
)

type state struct {
	elves set.Set[point.Point]
	order []string
}

func getBounds(s set.Set[point.Point]) (x0, y0, x1, y1 int) {
	s.Apply(func(e point.Point) {
		x0 = util.Min(x0, e.X)
		y0 = util.Min(y0, e.Y)
		x1 = util.Max(x1, e.X)
		y1 = util.Max(y1, e.Y)
	})
	return
}

type pp bool

func (p pp) String() string {
	if p {
		return "#"
	} else {
		return "."
	}
}

func printmap(s set.Set[point.Point]) string {
	minX, minY, maxX, maxY := getBounds(s)
	g, _ := grid.NewGrid[pp](minX, minY, maxX, maxY)
	s.Apply(func(e point.Point) { g.Set(e, true) })
	return g.String()
}

var (
	diag  = []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	check = map[string][]struct{ dx, dy int }{
		"N": []struct{ dx, dy int }{{-1, -1}, {0, -1}, {1, -1}},
		"S": []struct{ dx, dy int }{{-1, 1}, {0, 1}, {1, 1}},
		"W": []struct{ dx, dy int }{{-1, -1}, {-1, 0}, {-1, 1}},
		"E": []struct{ dx, dy int }{{1, -1}, {1, 0}, {1, 1}},
	}
	move = map[string]func(p point.Point) point.Point{
		"N": func(p point.Point) point.Point { return point.Point{p.X, p.Y - 1} },
		"S": func(p point.Point) point.Point { return point.Point{p.X, p.Y + 1} },
		"W": func(p point.Point) point.Point { return point.Point{p.X - 1, p.Y} },
		"E": func(p point.Point) point.Point { return point.Point{p.X + 1, p.Y} },
	}
)

func parseInput(r io.Reader) (out state) {
	out.elves = set.NewSet[point.Point]()
	out.order = []string{"N", "S", "W", "E"}
	y := 0
	util.Must(reader.LineReader(r, func(s string) error {
		for x, v := range strings.Split(s, "") {
			if v == "#" {
				out.elves.Add(point.Point{x, y})
			}
		}
		y++
		return nil
	}))
	return
}

func empty(elves set.Set[point.Point], e point.Point, d []struct{ dx, dy int }) bool {
	p := set.NewSet[point.Point]()
	for _, v := range d {
		p.Add(point.Point{e.X + v.dx, e.Y + v.dy})
	}
	return elves.Intersection(p).Len() == 0
}

func round(in set.Set[point.Point], order []string) (out set.Set[point.Point], done bool) {
	out = set.NewSet[point.Point]()
	proposed := make(map[point.Point]point.Point)
	count := make(map[point.Point]int)
	for e := range in {
		if empty(in, e, diag) {
			out.Add(e)
		} else {
			found := false
			for _, d := range order {
				if empty(in, e, check[d]) {
					next := move[d](e)
					proposed[e] = next
					count[next] = count[next] + 1
					found = true
					break
				}
			}
			if !found {
				out.Add(e)
			}
		}
	}
	fmt.Println(">>", out.Len(), len(proposed))
	if len(proposed) == 0 {
		done = true
	} else {
		for cur, next := range proposed {
			if count[next] > 1 {
				out.Add(cur)
			} else {
				out.Add(next)
			}
		}
	}
	return
}

func part1(input state) (result int) {
	done := false
	elves := input.elves.Copy()
	order := slices.Clone(input.order)
	for i := 0; i < 10 && !done; i++ {
		elves, done = round(elves, order)
		order[0], order[1], order[2], order[3] = order[1], order[2], order[3], order[0]
	}
	x0, y0, x1, y1 := getBounds(elves)
	return ((x1 - x0 + 1) * (y1 - y0 + 1)) - elves.Len()
}

func part2(input state) (result int) {
	elves := input.elves.Copy()
	order := slices.Clone(input.order)
	done := false
	for !done {
		elves, done = round(elves, order)
		order[0], order[1], order[2], order[3] = order[1], order[2], order[3], order[0]
		result++
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
