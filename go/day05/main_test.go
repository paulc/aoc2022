package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/paulc/aoc2022/util"
	"golang.org/x/exp/slices"
)

const testdata = `
    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2
`

var stacks, moves = parseInput(bytes.NewBufferString(strings.Trim(testdata, "\n")))

func TestPart1(t *testing.T) {
	result := part1(util.Map(stacks, func(s stack[string]) stack[string] { return slices.Clone(s) }), moves)
	if result != "CMZ" {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(util.Map(stacks, func(s stack[string]) stack[string] { return slices.Clone(s) }), moves)
	if result != "MCD" {
		t.Error(result)
	}
}
