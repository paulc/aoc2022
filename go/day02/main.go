package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
)

func parseInput(r io.Reader) array.Array[string] {
	return util.Must(array.ArrayReader[string](r, array.SplitWS, func(s string) (string, error) { return s, nil }))
}

var rules1 = map[string]map[string]int{
	"X": {"A": 4, "B": 1, "C": 7}, // Rock: 1 + result
	"Y": {"A": 8, "B": 5, "C": 2}, // Paper: 2 + result
	"Z": {"A": 3, "B": 9, "C": 6}, // Scissors: 3 + result
}

func part1(input array.Array[string]) (result int) {
	for _, v := range input {
		result += rules1[v[1]][v[0]]
	}
	return result
}

var rules2 = map[string]map[string]int{
	"X": {"A": 3, "B": 1, "C": 2}, // Lose: 0 + val
	"Y": {"A": 4, "B": 5, "C": 6}, // Draw: 3 + val
	"Z": {"A": 8, "B": 9, "C": 7}, // Win: 6 + val
}

func part2(input array.Array[string]) (result int) {
	for _, v := range input {
		result += rules2[v[1]][v[0]]
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
