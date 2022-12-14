package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/grid"
	"github.com/paulc/aoc2022/util/point"
	"github.com/paulc/aoc2022/util/reader"
)

type block int

const (
	empty block = iota
	rock
	sand
)

func (b block) String() string {
	switch b {
	case rock:
		return "#"
	case sand:
		return "o"
	default:
		return "."
	}
}

func parseInput(r io.Reader) (cave *grid.Grid[block]) {
	paths := [][]point.Point{}
	minX, maxX, maxY := math.MaxInt, 0, 0
	for _, v := range util.Must(reader.Lines(r)) {
		paths = append(paths, util.Map(strings.Split(v, " -> "), func(s string) point.Point {
			xy := util.Map(strings.Split(s, ","), func(s string) int { return util.Must(strconv.Atoi(s)) })
			maxX = util.Max(xy[0], maxX)
			minX = util.Min(xy[0], minX)
			maxY = util.Max(xy[1], maxY)
			return point.Point{xy[0], xy[1]}
		}))
	}
	cave = util.Must(grid.NewGrid[block](minX-maxY, 0, maxX+maxY, maxY+2))
	for _, p := range paths {
		for i := 0; i < len(p)-1; i++ {
			cave.DrawLine(p[i], p[i+1], rock)
		}
	}
	return
}

func dropSand(cave *grid.Grid[block], p, t point.Point) bool {
	for cave.CheckBounds(p) {
		found := false
		for _, next := range []point.Point{p.Move(0, 1), p.Move(-1, 1), p.Move(1, 1)} {
			if cave.Get(next) == empty {
				p, found = next, true
				break
			}
		}
		if !found {
			cave.Set(p, sand)
			return p != t
		}
	}
	return false
}

func part1(cave *grid.Grid[block]) (result int) {
	for ; dropSand(cave, point.Point{500, 0}, point.Point{-1, -1}); result++ {
	}
	return
}

func part2(cave *grid.Grid[block]) (result int) {
	cave.DrawLine(point.Point{cave.X0, cave.Y1}, point.Point{cave.X1, cave.Y1}, rock)
	for ; dropSand(cave, point.Point{500, 0}, point.Point{500, 0}); result++ {
	}
	return result + 1
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(util.Must(input.Copy())))
	fmt.Println("Part2:", part2(util.Must(input.Copy())))
}
