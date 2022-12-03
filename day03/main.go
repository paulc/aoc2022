package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
)

func priority(s string) int {
	if c := s[0]; c >= 'a' {
		return int(c - 'a' + 1)
	} else {
		return int(c - 'A' + 27)
	}
}

func parseInput(r io.Reader) (out [][]string) {
	util.Must(reader.LineReader(r, func(s string) error {
		out = append(out, strings.Split(s, ""))
		return nil
	}))
	return
}

func part1(input [][]string) (result int) {
	for _, v := range input {
		i, _ := set.NewSetFrom(v[:len(v)/2]).Intersection(set.NewSetFrom(v[len(v)/2:])).Pop()
		result += priority(i)
	}
	return result
}

func part2(input [][]string) (result int) {
	for _, v := range util.Take(input, 3) {
		badge, _ := set.NewSetFrom(v[0]).Intersection(set.NewSetFrom(v[1])).Intersection(set.NewSetFrom(v[2])).Pop()
		result += priority(badge)
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
