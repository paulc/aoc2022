package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/paulc/aoc2022/util/must"
	"github.com/paulc/aoc2022/util/reader"
	"golang.org/x/exp/slices"
)

func parseInput(r io.Reader) []int {
	elves := []int{0}
	i := 0
	reader.LineReader(r, func(s string) error {
		if s == "" {
			elves = append(elves, 0)
			i++
		} else {
			elves[i] += must.Must(strconv.Atoi(s))
		}
		return nil
	})
	slices.SortFunc(elves, func(a, b int) bool { return a < b })
	return elves
}

func part1(elves []int) int {
	return elves[len(elves)-1]
}

func part2(elves []int) int {
	return elves[len(elves)-1] + elves[len(elves)-2] + elves[len(elves)-3]
}

func main() {
	input := parseInput(must.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
