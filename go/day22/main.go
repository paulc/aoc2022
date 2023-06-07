package main

import (
	"fmt"
	"os"

	"github.com/paulc/aoc2022/util"
)

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input, part2_map))
}
