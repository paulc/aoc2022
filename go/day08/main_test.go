package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
30373
25512
65332
33549
35390
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 21 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 8 {
		t.Error(result)
	}
}
