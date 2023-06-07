package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"golang.org/x/exp/slices"
)

type stack[T any] []T

func (s *stack[T]) Pop(n int) (out []T) {
	out = (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return
}

func (s *stack[T]) Push(v []T) {
	*s = append(*s, v...)
}

func parseInput(r io.Reader) (stacks []stack[string], moves [][]int) {
	head, tail := util.Split(util.Must(reader.Lines(r)), func(s string) bool { return s == "" })
	moves = util.Map(tail, func(s string) []int { return util.Must(util.SlurpInt(s)) })
	n := (len(head[0]) + 1) / 4
	stacks = make([]stack[string], n)
	for i := len(head) - 2; i >= 0; i-- {
		for j := 0; j < n; j++ {
			c := head[i][j*4+1]
			if c != ' ' {
				stacks[j] = append(stacks[j], string(c))
			}
		}
	}
	return
}

func part1(stacks []stack[string], moves [][]int) (result string) {
	for _, v := range moves {
		for i := 0; i < v[0]; i++ {
			stacks[v[2]-1].Push(stacks[v[1]-1].Pop(1))
		}
	}
	return strings.Join(util.Map(stacks, func(s stack[string]) string { return s[len(s)-1] }), "")
}

func part2(stacks []stack[string], moves [][]int) (result string) {
	for _, v := range moves {
		stacks[v[2]-1].Push(stacks[v[1]-1].Pop(v[0]))
	}
	return strings.Join(util.Map(stacks, func(s stack[string]) string { return s[len(s)-1] }), "")
}

func main() {
	stacks, moves := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(util.Map(stacks, func(s stack[string]) stack[string] { return slices.Clone(s) }), moves))
	fmt.Println("Part2:", part2(util.Map(stacks, func(s stack[string]) stack[string] { return slices.Clone(s) }), moves))
}
