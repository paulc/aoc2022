package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
)

type Tree struct {
	height  int
	visible bool
}

func parseInput(r io.Reader) array.Array[Tree] {
	return util.Must(array.ArrayReader(r, array.MakeStringSplitter(""), func(s string) (Tree, error) {
		return Tree{height: util.Must(strconv.Atoi(s))}, nil
	}))
}

func part1(input array.Array[Tree]) (result int) {
	h, w := len(input), len(input[0])
	for y := 0; y < h; y++ {
		for _, v := range []struct{ start, dx int }{{0, 1}, {w - 1, -1}} {
			for x, max := v.start, -1; x >= 0 && x < w; x += v.dx {
				if input[y][x].height > max {
					if !input[y][x].visible {
						input[y][x].visible = true
						result++
					}
					max = input[y][x].height
				}
			}
		}
	}
	for x := 0; x < w; x++ {
		for _, v := range []struct{ start, dy int }{{0, 1}, {h - 1, -1}} {
			for y, max := v.start, -1; y >= 0 && y < h; y += v.dy {
				if input[y][x].height > max {
					if !input[y][x].visible {
						input[y][x].visible = true
						result++
					}
					max = input[y][x].height
				}
			}
		}
	}
	return result
}

func part2(input array.Array[Tree]) (result int) {
	h := len(input)
	w := len(input[0])
	input.Each(func(e array.ArrayElement[Tree]) {
		score := 1
		for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
			view := 0
			for x, y := e.X+v.dx, e.Y+v.dy; x >= 0 && x < w && y >= 0 && y < h; x, y = x+v.dx, y+v.dy {
				view++
				if input[y][x].height >= e.Val.height {
					break
				}
			}
			score *= view
		}
		if score > result {
			result = score
		}
	})
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
