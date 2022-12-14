package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/paulc/aoc2022/util"
)

const testdata = `
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(util.Must(input.Copy()))
	if result != 24 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(util.Must(input.Copy()))
	if result != 93 {
		t.Error(result)
	}
}
