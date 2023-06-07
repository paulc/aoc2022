package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type puzzle struct {
}

func parseInput(r io.Reader) (out puzzle) {
	util.Must(reader.LineReader(r, func(s string) error {
		return nil
	}))
	return
}

func part1(input puzzle) (result int) {
	return result
}

func part2(input puzzle) (result int) {
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
