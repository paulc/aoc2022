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

type board struct {
	rows  [][7]bool
	top   int
	start int
}

type rock struct {
	points []offset
	height int
}

var rocks = []rock{
	{[]offset{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, 1},
	{[]offset{{1, 0}, {0, -1}, {1, -1}, {2, -1}, {1, -2}}, 3},
	{[]offset{{0, -2}, {1, -2}, {2, -2}, {2, -1}, {2, 0}}, 3},
	{[]offset{{0, 0}, {0, -1}, {0, -2}, {0, -3}}, 4},
	{[]offset{{0, 0}, {0, -1}, {1, 0}, {1, -1}}, 2},
}

var jetMove = map[byte]offset{'<': {-1, 0}, '>': {1, 0}}

func (b *board) AddRock(r rock, jet <-chan byte) {
	top := b.top + 3 + (r.height - 1)
	for top > len(b.rows)-1 {
		b.rows = append(b.rows, [7]bool{})
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

func (b *board) Draw(pos xy, r rock) {
	for _, v := range r.points {
		p := pos.Move(v)
		b.rows[p.y][p.x] = true
	}
}

func (b *board) Check(pos xy, r rock) bool {
	for _, v := range r.points {
		p := pos.Move(v)
		if p.y < 0 || p.x < 0 || p.x >= 7 || b.rows[p.y][p.x] {
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

func cycle[T any](in []T) <-chan T {
	out := make(chan T)
	go func() {
		for {
			for _, v := range in {
				out <- v
			}
		}
	}()
	return out
}

func parseInput(r io.Reader) (out startData) {
	return bytes.TrimSpace(util.Must(io.ReadAll(r)))
}

func part1(input startData) (result int) {
	b := board{}
	jets := cycle(input)
	i := 0
	for v := range cycle(rocks) {
		b.AddRock(v, jets)
		i++
		if i >= 2022 {
			break
		}
	}
	return b.top
}

func part2(input startData) (result int) {
	return
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
