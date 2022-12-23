package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/array"
	"github.com/paulc/aoc2022/util/reader"
)

type tile int8

const (
	void tile = iota
	open
	solid
)

type facing int8

const (
	R facing = iota
	D
	L
	U
	nFacing
)

func (f facing) String() string {
	return map[facing]string{R: "R", D: "D", L: "L", U: "U"}[f]
}

var dxy = map[facing]struct{ dx, dy int }{
	R: {1, 0}, L: {-1, 0}, U: {0, -1}, D: {0, 1},
}

func strToTile(s string) (tile, error) {
	switch s {
	case ".":
		return open, nil
	case "#":
		return solid, nil
	case " ":
		return void, nil
	default:
		return void, errors.New("Invalid tile")
	}
}

func (t tile) String() string {
	switch t {
	case open:
		return "."
	case solid:
		return "#"
	case void:
		return " "
	default:
		return "?"
	}
}

type move struct {
	turn      bool
	count     int
	direction string
}

type pos struct {
	x, y int
	dir  facing
}

type startData struct {
	cave           array.Array[tile]
	moves          []move
	start, current pos
	w, h           int
}

/*
func (s startData) String() string {
	return fmt.Sprintf("%s\n\nw=%d h=%d c=(%d,%d) %s\n", s.cave.Copy().Set(s.current.x, s.current.y, 99).String(), s.w, s.h, s.current.x, s.current.y, s.current.dir)
}
*/

func (s *startData) Move(m move) {
	if m.turn {
		if m.direction == "R" {
			s.current.dir = (s.current.dir + 1) % nFacing
		} else {
			s.current.dir = (s.current.dir + 3) % nFacing
		}
	} else {
		delta := dxy[s.current.dir]
		for i := 0; i < m.count; i++ {
			onmap := false
			x, y := s.current.x, s.current.y
			for !onmap {
				x1, y1 := x+delta.dx, y+delta.dy
				if x1 < 0 {
					x1 = s.w - 1
				}
				if x1 >= s.w {
					x1 = 0
				}
				if y1 < 0 {
					y1 = s.h - 1
				}
				if y1 >= s.h {
					y1 = 0
				}
				switch s.cave[y1][x1] {
				case solid:
					return // Blocked
				case open:
					s.current.x, s.current.y = x1, y1
					onmap = true
				case void:
					x, y = x1, y1
				}
			}
		}
	}
	return
}

func parseInput(r io.Reader) (out startData) {
	head, tail, _ := reader.HeadFunc(r, func(b []byte) bool { return bytes.Equal(b, []byte{'\n'}) }, true)
	out.cave = util.Must(array.ArrayReader[tile](&head, array.MakeStringSplitter(""), strToTile))
	re := regexp.MustCompile(`(\d+|[LR])`)
	b, i := tail.Bytes(), 0
	for m := re.FindIndex(b[i:]); m != nil; {
		a := string(b[i:][m[0]:m[1]])
		mv := move{}
		if a == "L" || a == "R" {
			mv.direction = a
			mv.turn = true
		} else {
			mv.count = util.Must(strconv.Atoi(a))
		}
		out.moves = append(out.moves, mv)
		i += m[1]
		m = re.FindIndex(b[i:])
	}
	for out.start.x = 0; out.cave[0][out.start.x] != open; out.start.x++ {
	}
	out.current = out.start
	out.h = len(out.cave)
	for i := 0; i < out.h; i++ {
		if l := len(out.cave[i]); l > out.w {
			out.w = l
		}
	}
	for i := 0; i < out.h; i++ {
		if l := len(out.cave[i]); l < out.w {
			for j := 0; j < out.w-l; j++ {
				out.cave[i] = append(out.cave[i], void)
			}
		}
	}
	return
}

func part1(input startData) (result int) {
	fmt.Println(input.cave)
	for _, v := range input.moves {
		input.Move(v)
	}
	return (input.current.y+1)*1000 + (input.current.x+1)*4 + int(input.current.dir)
}

func part2(input startData) (result int) {
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}