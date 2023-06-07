package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/set"
)

func parseInput(r io.Reader) []byte {
	return util.Must(io.ReadAll(r))
}

func findStart(input []byte, startLen int) (n int) {
	for i := 0; i < len(input)-startLen; {
		unique := set.NewSetFrom(input[i : i+startLen]).Len()
		if unique == startLen {
			return i + startLen
		}
		i += startLen - unique
	}
	return
}

func part1(input []byte) (result int) {
	return findStart(input, 4)
}

func part2(input []byte) (result int) {
	return findStart(input, 14)
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
