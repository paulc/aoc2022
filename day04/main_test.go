package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 2 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 4 {
		t.Error(result)
	}
}
