package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/paulc/aoc2022/util/array"
	"github.com/paulc/aoc2022/util/must"
)

func parseInput(r io.Reader) array.Array[int] {
	return must.Must(array.ArrayReader(r, array.SplitWS, func(s string) (int, error) { return strconv.Atoi(s) }))
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
	input := parseInput(must.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
