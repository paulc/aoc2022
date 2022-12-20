package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type num struct {
	val, startPos int
}

type startData struct {
	data []num
	pos  map[int]int
}

func (s *startData) wrap(n int) int {
	n = n % len(s.data)
	if n < 0 {
		n = len(s.data) + n
	}
	return n
}

func (s *startData) move(pos int) {
	val := s.data[pos].val
	if val > 0 {
		for i := pos; i < pos+val; i++ {
			s.data[s.wrap(i)], s.data[s.wrap(i+1)] = s.data[s.wrap(i+1)], s.data[s.wrap(i)]
			s.pos[s.data[s.wrap(i)].startPos] = s.wrap(i)
			s.pos[s.data[s.wrap(i+1)].startPos] = s.wrap(i + 1)
		}
	} else {
		for i := pos; i > pos+val; i-- {
			s.data[s.wrap(i)], s.data[s.wrap(i-1)] = s.data[s.wrap(i-1)], s.data[s.wrap(i)]
			s.pos[s.data[s.wrap(i)].startPos] = s.wrap(i)
			s.pos[s.data[s.wrap(i-1)].startPos] = s.wrap(i - 1)
		}
	}
}

func parseInput(r io.Reader) (out startData) {
	out.pos = make(map[int]int)
	i := 0
	util.Must(reader.LineReader(r, func(s string) error {
		out.data = append(out.data, num{util.Must(strconv.Atoi(s)), i})
		out.pos[i] = i
		i++
		return nil
	}))
	return
}

func part1(input startData) (result int) {
	for i := 0; i < len(input.data); i++ {
		input.move(input.pos[i])
	}
	for i, v := range input.data {
		if v.val == 0 {
			result = input.data[input.wrap(i+1000)].val + input.data[input.wrap(i+2000)].val + input.data[input.wrap(i+3000)].val
			break
		}
	}
	return
}

func part2(input startData) (result int) {
	return
}

func main() {
	fmt.Println("Part1:", part1(parseInput(util.Must(os.Open("input")))))
	fmt.Println("Part2:", part2(parseInput(util.Must(os.Open("input")))))
}
