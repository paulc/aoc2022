package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
)

var moves = []xyz{{1, 0, 0}, {-1, 0, 0}, {0, 1, 0}, {0, -1, 0}, {0, 0, 1}, {0, 0, -1}}

type xyz struct{ x, y, z int }

func (p xyz) adjacent() []xyz {
	out := make([]xyz, 6)
	for i, v := range moves {
		out[i] = xyz{p.x + v.x, p.y + v.y, p.z + v.z}
	}
	return out
}

func (p xyz) inside(p1, p2 xyz) bool {
	return p.x >= p1.x && p.x <= p2.x && p.y >= p1.y && p.y <= p2.y && p.z >= p1.z && p.z <= p2.z
}

func parseInput(r io.Reader) (out set.Set[xyz]) {
	out = set.NewSetFrom(util.Map(util.Must(reader.Lines(r)), func(s string) (p xyz) {
		util.Must(fmt.Sscanf(s, "%d,%d,%d", &p.x, &p.y, &p.z))
		return
	}))
	return
}

func checkEnclosed(cubes set.Set[xyz], p, pmin, pmax xyz) bool {
	blocked := make([]bool, 6)
	for i, v := range moves {
		for p2 := p; p2.inside(pmin, pmax); {
			p2 = xyz{p2.x + v.x, p2.y + v.y, p2.z + v.z}
			if cubes.Has(p2) {
				blocked[i] = true
				break
			}
		}
	}
	return util.Reduce(blocked, func(a, b bool) bool { return a && b }, true)
}

func part1(input set.Set[xyz]) (result int) {
	for p := range input {
		for _, a := range p.adjacent() {
			if !input.Has(a) {
				result++
			}
		}
	}
	return result
}

func part2(input set.Set[xyz]) (result int) {
	var pmin, pmax xyz
	for p := range input {
		pmin = xyz{util.Min(pmin.x, p.x-1), util.Min(pmin.y, p.y-1), util.Min(pmin.z, p.z-1)}
		pmax = xyz{util.Max(pmax.x, p.x+1), util.Max(pmax.y, p.y+1), util.Max(pmax.z, p.z+1)}
	}
	for x := pmin.x; x < pmax.x; x++ {
		for y := pmin.y + 1; y < pmax.y; y++ {
			for z := pmin.z; z < pmax.z; z++ {
				if !input.Has(xyz{x, y, z}) {
					if checkEnclosed(input, xyz{x, y, z}, pmin, pmax) {
						input.Add(xyz{x, y, z})
					}
				}

			}
		}
	}
	for p := range input {
		for _, a := range p.adjacent() {
			if !input.Has(a) {
				result++
			}
		}
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
