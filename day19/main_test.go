package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 33 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 3472 {
		t.Error(result)
	}
}
