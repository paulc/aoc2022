package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type puzzle []int

func fromSnafu(s string) (result int) {
	for i := 0; i < len(s); i++ {
		result = result*5 + map[byte]int{'2': 2, '1': 1, '0': 0, '-': -1, '=': -2}[s[i]]
	}
	return
}

func toSnafu(i int) string {
	result := []string{}
	for r := 0; i > 0; {
		i, r = i/5, i%5
		result = append([]string{map[int]string{0: "0", 1: "1", 2: "2", 3: "=", 4: "-"}[r]}, result...)
		if r > 2 {
			i += 1
		}
	}
	return strings.Join(result, "")
}

func parseInput(r io.Reader) (out puzzle) {
	util.Must(reader.LineReader(r, func(s string) error {
		out = append(out, fromSnafu(s))
		return nil
	}))
	return
}

func part1(input puzzle) (result string) {
	sum := 0
	util.Apply(input, func(i int) { sum += i })
	return toSnafu(sum)
}

func part2(input puzzle) (result string) {
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
