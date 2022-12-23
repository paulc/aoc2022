package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/paulc/aoc2022/util"
)

const testdata = `
....#..
..###.#
#...#.#
.#...##
#.###..
##.#.##
.#..#..
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 110 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 20 {
		t.Error(result)
	}
}

func BenchPart1(t *testing.T) {
	result := part1(parseInput(util.Must(os.Open("input"))))
	if result != 4208 {
		t.Error(result)
	}
}
