package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type startData struct {
}

func parseInput(r io.Reader) (out startData) {
	util.Must(reader.LineReader(r, func(s string) error {
	}))
	return
}

func part1(input startData) (result int) {
	return result
}

func part2(input startData) (result int) {
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
