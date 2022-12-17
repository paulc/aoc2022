package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 3068 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 1514285714288 {
		t.Error(result)
	}
}
