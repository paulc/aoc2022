package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
)

type startData []byte

type xy struct{ x, y int }

func (p xy) Move(o offset) xy {
	return xy{p.x + o.dx, p.y + o.dy}
}

type offset struct{ dx, dy int }

const boardRows = 100

type signature struct {
	rows [boardRows / 4]int8
	rock int
}

type board struct {
	rows  [boardRows][7]bool
	top   int
	start int
	cache map[signature][2]int
}

type rock struct {
	points []offset
	height int
	index  int
}

var rocks = []rock{
	{[]offset{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 1, 0},
	{[]offset{{1, 0}, {0, -1}, {1, -1}, {2, -1}, {1, -2}}, 3, 1},
	{[]offset{{0, -2}, {1, -2}, {2, -2}, {2, -1}, {2, 0}}, 3, 2},
	{[]offset{{0, 0}, {0, -1}, {0, -2}, {0, -3}}, 4, 3},
	{[]offset{{0, 0}, {0, -1}, {1, 0}, {1, -1}}, 2, 4},
}

var jetMove = map[byte]offset{'<': {-1, 0}, '>': {1, 0}}

func toInt8(in [7]bool) (out int8) {
	for i := 0; i < 7; i++ {
		out = out << 1
		if in[i] {
			out += 1
		}
	}
	return
}

func (b *board) AddRock(r rock, jet <-chan byte) {

	top := b.top + 3 + (r.height - 1)

	if top-b.start >= boardRows {
		for i := boardRows / 2; i < boardRows; i++ {
			b.rows[i-boardRows/2] = b.rows[i]
			b.rows[i] = [7]bool{}
		}
		b.start += boardRows / 2
	}
	pos := xy{2, top}

	for {
		next := pos.Move(jetMove[<-jet])
		if b.Check(next, r) {
			pos = next
		}
		next = pos.Move(offset{0, -1})
		if b.Check(next, r) {
			pos = next
		} else {
			break
		}
	}
	if pos.y+1 > b.top {
		b.top = pos.y + 1
	}
	b.Draw(pos, r)
}

func (b *board) CalculateSignature(r rock) (sig signature) {
	sig.rock = r.index
	for i := 0; i < boardRows/4; i++ {
		sig.rows[i] = toInt8(b.rows[b.top-b.start-i])
	}
	return
}

func (b *board) Draw(pos xy, r rock) {
	for _, v := range r.points {
		p := pos.Move(v)
		b.rows[p.y-b.start][p.x] = true
	}
}

func (b *board) Check(pos xy, r rock) bool {
	for _, v := range r.points {
		p := pos.Move(v)
		if p.y < 0 || p.x < 0 || p.x >= 7 || b.rows[p.y-b.start][p.x] {
			return false
		}
	}
	return true
}

func (b board) String() string {
	h := len(b.rows)
	out := make([]string, h)
	for i := 0; i < h; i++ {
		out[i] = fmt.Sprintf("%4d %s", h-i-1, strings.Join(util.Map(b.rows[h-i-1][:], func(b bool) string {
			if b {
				return "#"
			} else {
				return "."
			}
		}), ""))
	}
	out = append(out, fmt.Sprintf("Top: %d\n\n", b.top))
	return strings.Join(out, "\n")
}

func parseInput(r io.Reader) (out startData) {
	return bytes.TrimSpace(util.Must(io.ReadAll(r)))
}

func run(input startData, nrocks int) int {
	b := board{cache: make(map[signature][2]int)}
	jets := util.Cycle(input)
	i, found := 0, false
	for v := range util.Cycle(rocks) {
		b.AddRock(v, jets)
		i++
		if b.top > boardRows/2 && !found {
			sig := b.CalculateSignature(v)
			if prev, ok := b.cache[sig]; ok {
				found = true
				cycleLen := i - prev[0]
				cycleHeight := b.top - prev[1]
				ncycles := (nrocks - i) / cycleLen
				i += ncycles * cycleLen
				b.top += ncycles * cycleHeight
				b.start += ncycles * cycleHeight
			} else {
				b.cache[sig] = [2]int{i, b.top}
			}
		}
		if i == nrocks {
			break
		}
	}
	return b.top
}

func part1(input startData) (result int) {
	return run(input, 2022)
}

func part2(input startData) (result int) {
	return run(input, 1000000000000)
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
