package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
)

type Move struct {
	dir   string
	count int
}

type XY struct{ x, y int }

func (p XY) Move(d XY) XY {
	return XY{p.x + d.x, p.y + d.y}
}

func (p XY) Offset(p1 XY) XY {
	return XY{p1.x - p.x, p1.y - p.y}
}

var (
	moveHead = map[string]XY{
		"L": XY{-1, 0}, "R": XY{1, 0}, "U": XY{0, 1}, "D": XY{0, -1},
	}
	moveTail = map[XY]XY{
		XY{-2, 2}: XY{1, -1}, XY{-1, 2}: XY{1, -1}, XY{0, 2}: XY{0, -1}, XY{1, 2}: XY{-1, -1}, XY{2, 2}: XY{-1, -1},
		XY{-2, -2}: XY{1, 1}, XY{-1, -2}: XY{1, 1}, XY{0, -2}: XY{0, 1}, XY{1, -2}: XY{-1, 1}, XY{2, -2}: XY{-1, 1},
		XY{-2, 1}: XY{1, -1}, XY{-2, 0}: XY{1, 0}, XY{-2, -1}: XY{1, 1},
		XY{2, 1}: XY{-1, -1}, XY{2, 0}: XY{-1, 0}, XY{2, -1}: XY{-1, 1},
	}
)

func moveRope(n int, moves []Move) int {
	rope := make([]XY, n)
	visited := set.NewSetFrom([]XY{XY{0, 0}})
	for _, v := range moves {
		for i := 0; i < v.count; i++ {
			rope[0] = rope[0].Move(moveHead[v.dir])
			for j := 1; j < n; j++ {
				rope[j] = rope[j].Move(moveTail[rope[j-1].Offset(rope[j])])
			}
			visited.Add(rope[n-1])
		}
	}
	return visited.Len()
}

func parseInput(r io.Reader) (out []Move) {
	return util.Map(util.Must(reader.Lines(r)), func(s string) (m Move) {
		util.Must(fmt.Sscanf(s, "%s %d", &m.dir, &m.count))
		return
	})
}

func part1(input []Move) (result int) {
	return moveRope(2, input)
}

func part2(input []Move) (result int) {
	return moveRope(10, input)
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
