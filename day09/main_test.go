package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata1 = `
R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2
`

var testdata2 = `
R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20
`

var (
	input1 = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata1)))
	input2 = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata2)))
)

func TestPart1(t *testing.T) {
	result := part1(input1)
	if result != 13 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input1)
	if result != 1 {
		t.Error(result)
	}
}

func TestPart2a(t *testing.T) {
	result := part2(input2)
	if result != 36 {
		t.Error(result)
	}
}
