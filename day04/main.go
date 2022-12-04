package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

func parseInput(r io.Reader) (out [][]int) {
	util.Must(reader.LineReader(r, func(s string) error {
		out = append(out, util.Must(util.SlurpInt(s)))
		return nil
	}))
	return out
}

func part1(input [][]int) (result int) {
	for _, v := range input {
		result += map[bool]int{true: 1, false: 0}[(v[0] <= v[2] && v[1] >= v[3]) || (v[2] <= v[0] && v[3] >= v[1])]
	}
	return result
}

func part2(input [][]int) (result int) {
	for _, v := range input {
		result += map[bool]int{true: 1, false: 0}[(v[0] <= v[3] && v[1] >= v[2]) || (v[2] <= v[1] && v[3] >= v[0])]
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
