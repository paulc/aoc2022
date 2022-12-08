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
	h := len(input)
	w := len(input[0])
	for y := 0; y < h; y++ {
		for x, max := 0, -1; x < w; x++ {
			if input[y][x].height > max {
				input[y][x].visible = true
				max = input[y][x].height
			}
		}
		for x, max := w-1, -1; x >= 0; x-- {
			if input[y][x].height > max {
				input[y][x].visible = true
				max = input[y][x].height
			}
		}
	}
	for x := 0; x < w; x++ {
		for y, max := 0, -1; y < h; y++ {
			if input[y][x].height > max {
				input[y][x].visible = true
				max = input[y][x].height
			}
		}
		for y, max := h-1, -1; y >= 0; y-- {
			if input[y][x].height > max {
				input[y][x].visible = true
				max = input[y][x].height
			}
		}
	}
	count := 0
	input.Each(func(e array.ArrayElement[Tree]) {
		if e.Val.visible {
			count++
		}
	})
	return count
}

func part2(input array.Array[Tree]) (result int) {
	h := len(input)
	w := len(input[0])
	input.Each(func(e array.ArrayElement[Tree]) {
		score := 1
		for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
			x, y, view := e.X, e.Y, 0
			for {
				x, y = x+v.dx, y+v.dy
				if x >= 0 && x < w && y >= 0 && y < h {
					view++
					if input[y][x].height >= e.Val.height {
						break
					}
				} else {
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
