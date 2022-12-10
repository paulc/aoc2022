package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
)

type Crt [240]bool

func (crt Crt) String() string {
	x := util.Map(crt[:], func(p bool) byte { return map[bool]byte{true: '#', false: '.'}[p] })
	out := []string{}
	for i := 0; i < len(crt); i += 40 {
		out = append(out, string(x[i:i+40]))
	}
	return strings.Join(out, "\n")
}

func parseInput(r io.Reader) [][]string {
	return util.Map(util.Must(reader.Lines(r)), func(s string) []string { return strings.Split(s, " ") })
}

func runCpu(input [][]string, f func(cycle, X int)) {
	X, cycle := 1, 0
	for _, v := range input {
		if v[0] == "noop" {
			cycle += 1
			f(cycle, X)
		} else if v[0] == "addx" {
			for i := 0; i < 2; i++ {
				cycle += 1
				f(cycle, X)
			}
			X += util.Must(strconv.Atoi(v[1]))
		}
	}
}

func part1(input [][]string) (result int) {
	interesting := set.NewSetFrom([]int{20, 60, 100, 140, 180, 220})
	values := []int{}
	runCpu(input, func(cycle, X int) {
		if interesting.Has(cycle) {
			values = append(values, cycle*X)
		}
	})
	return util.Reduce(values, func(a, b int) int { return a + b }, 0)
}

func part2(input [][]string) (result string) {
	crt := Crt{}
	runCpu(input, func(cycle, X int) {
		if pos := cycle - 1; (pos%40) >= X-1 && (pos%40) <= X+1 {
			crt[pos] = true
		}
	})
	return crt.String()
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:")
	fmt.Println(part2(input))
}
