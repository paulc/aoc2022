package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"golang.org/x/exp/slices"
)

type Monkey struct {
	items     []int
	op        func(int) int
	test      func(int) bool
	next      map[bool]int
	inspected int
}

func makeOp(op []string) func(int) int {
	return func(v int) int {
		var a, b int
		if op[2] == "old" {
			a = v
		} else {
			a = util.Must(strconv.Atoi(op[2]))
		}
		if op[4] == "old" {
			b = v
		} else {
			b = util.Must(strconv.Atoi(op[4]))
		}
		if op[3] == "+" {
			return a + b
		} else {
			return a * b
		}
	}
}

func parseInput(r io.Reader) (out []Monkey, lcm int) {
	m := Monkey{next: make(map[bool]int)}
	re := regexp.MustCompile(": ?")
	lcm = 1
	for _, v := range util.Must(reader.Lines(r)) {
		if v == "" {
			out = append(out, m)
			m = Monkey{next: make(map[bool]int)}
		} else {
			p := re.Split(strings.TrimSpace(v), 2)
			switch p[0] {
			case "Starting items":
				m.items = util.Map(strings.Split(p[1], ", "), func(s string) int { return util.Must(strconv.Atoi(s)) })
			case "Operation":
				m.op = makeOp(strings.Fields(p[1]))
			case "Test":
				div := util.Must(strconv.Atoi(strings.Fields(p[1])[2]))
				lcm = lcm * div
				m.test = func(i int) bool { return i%div == 0 }
			case "If true":
				m.next[true] = util.Must(strconv.Atoi(strings.Fields(p[1])[3]))
			case "If false":
				m.next[false] = util.Must(strconv.Atoi(strings.Fields(p[1])[3]))
			}
		}
	}
	out = append(out, m)
	return
}

func shuffle(monkeys []Monkey, div int, lcm int) {
	for i := range monkeys {
		for _, item := range monkeys[i].items {
			monkeys[i].inspected++
			wl := (monkeys[i].op(item) / div) % lcm
			next := monkeys[i].next[monkeys[i].test(wl)]
			monkeys[next].items = append(monkeys[next].items, wl)
		}
		monkeys[i].items = []int{}
	}
}

func part1(monkeys []Monkey, lcm int) (result int) {
	for round := 0; round < 20; round++ {
		shuffle(monkeys, 3, lcm)
	}
	slices.SortFunc(monkeys, func(a, b Monkey) bool { return a.inspected > b.inspected })
	return monkeys[0].inspected * monkeys[1].inspected
}

func part2(monkeys []Monkey, lcm int) (result int) {
	for round := 0; round < 10000; round++ {
		shuffle(monkeys, 1, lcm)
	}
	slices.SortFunc(monkeys, func(a, b Monkey) bool { return a.inspected > b.inspected })
	return monkeys[0].inspected * monkeys[1].inspected
}

func main() {
	fmt.Println("Part1:", part1(parseInput(util.Must(os.Open("input")))))
	fmt.Println("Part2:", part2(parseInput(util.Must(os.Open("input")))))
}
