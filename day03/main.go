package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util/array"
)

func parseInput(r io.Reader) array.Array[string] {
	a, err := array.ArrayReader[string](r, array.SplitWS, func(s string) (string, error) { return s, nil })
	if err != nil {
		panic(err)
	}
	return a
}

func part1(input array.Array[string]) int {
	result := 0
	return result
}

func part2(input array.Array[string]) int {
	result := 0
	return result
}

func main() {
	r, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	input := parseInput(r)
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
