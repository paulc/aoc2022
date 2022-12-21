package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type node struct {
	d1, d2 string
	op     func(a, b int) int
	val    int
	hasVal bool
}

var fn = map[string]func(a, b int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

type startData map[string]*node

func (s startData) Value(id string) int {
	n := s[id]
	if !n.hasVal {
		return n.op(s.Value(n.d1), s.Value(n.d2))
	} else {
		return n.val
	}
}

func parseInput(r io.Reader) (out startData) {
	out = make(startData)
	util.Must(reader.LineReader(r, func(s string) error {
		n := node{}
		s1 := strings.Split(s, ": ")
		s2 := strings.Split(s1[1], " ")
		if len(s2) == 1 {
			n.val = util.Must(strconv.Atoi(s2[0]))
			n.hasVal = true
		} else {
			n.d1, n.op, n.d2 = s2[0], fn[s2[1]], s2[2]
		}
		out[s1[0]] = &n
		return nil
	}))
	return
}

func guess(f func(i int) int, target int) (guess1, guess2 int) {
	slope := float64(f(1000)-f(0)) / 1000
	guess1 = int(float64(target-f(0)) / slope)
	diff := f(guess1) - target
	guess2 = guess1 - int(float64(diff*2)/slope)
	return
}

func bsearch(f func(i int) int, start, end, target int) int {
	mid := (start + end) / 2
	vmid, vend := f(mid), f(end)
	if vmid == target {
		return mid
	}
	if (vmid > target && vend > target) || (vmid < target && vend < target) {
		return bsearch(f, start, mid, target)
	} else {
		return bsearch(f, mid, end, target)
	}
}

func part1(input startData) (result int) {
	return input.Value("root")
}

func part2(input startData) (result int) {
	var f func(i int) int
	d1 := func(i int) int { input["humn"].val = i; return input.Value(input["root"].d1) }
	d2 := func(i int) int { input["humn"].val = i; return input.Value(input["root"].d2) }
	target := 0
	if d1(9999) != d1(-9999) {
		f, target = d1, d2(0)
	} else {
		f, target = d2, d1(0)
	}
	g1, g2 := guess(f, target)
	// Can have multiple 'correct' answers - pick lowest
	for x := bsearch(f, g1, g2, target); d1(x) == d2(x); x-- {
		result = x
	}
	return
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
