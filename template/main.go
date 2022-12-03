package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
)

func parseInput(r io.Reader) array.Array[int] {
	return util.Must(array.ArrayReader(r, array.SplitWS, func(s string) (int, error) { return strconv.Atoi(s) }))
}

func part1(input array.Array[int]) int {
	result := 0
	return result
}

func part2(input array.Array[int]) int {
	result := 0
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
