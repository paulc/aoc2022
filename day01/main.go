package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

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
			cal, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			elves[i] += cal
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
	r, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	input := parseInput(r)
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
