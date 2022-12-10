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

type Cpu struct {
	X int
}

func parseInput(r io.Reader) [][]string {
	return util.Map(util.Must(reader.Lines(r)), func(s string) []string { return strings.Split(s, " ") })
}

func part1(input [][]string) (result int) {
	cpu := Cpu{X: 1}
	cycle := 0
	interesting := set.NewSetFrom([]int{20, 60, 100, 140, 180, 220})
	values := []int{}
	for _, v := range input {
		if v[0] == "noop" {
			cycle += 1
			if interesting.Has(cycle) {
				values = append(values, cycle*cpu.X)
			}
		} else if v[0] == "addx" {
			for i := 0; i < 2; i++ {
				cycle += 1
				if interesting.Has(cycle) {
					values = append(values, cycle*cpu.X)
				}
			}
			cpu.X += util.Must(strconv.Atoi(v[1]))
		}
	}
	return util.Reduce(values, func(a, b int) int { return a + b }, 0)
}

type Crt [240]bool

func (crt Crt) String() string {
	out := []string{}
	for r := 0; r < 6; r++ {
		row := [40]byte{}
		for c := 0; c < 40; c++ {
			if crt[(r*40)+c] {
				row[c] = '#'
			} else {
				row[c] = ' '
			}
		}
		out = append(out, string(row[:]))
	}
	return strings.Join(out, "\n")
}

func part2(input [][]string) (result string) {
	cpu := Cpu{X: 1}
	cycle := 0
	crt := Crt{}
	for _, v := range input {
		if v[0] == "noop" {
			cycle += 1
			pos := cycle - 1
			if (pos%40) >= cpu.X-1 && (pos%40) <= cpu.X+1 {
				crt[pos] = true
			}
		} else if v[0] == "addx" {
			for i := 0; i < 2; i++ {
				cycle += 1
				pos := cycle - 1
				if (pos%40) >= cpu.X-1 && (pos%40) <= cpu.X+1 {
					crt[pos] = true
				}
			}
			cpu.X += util.Must(strconv.Atoi(v[1]))
		}
	}
	return crt.String()
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:")
	fmt.Println(part2(input))
}
