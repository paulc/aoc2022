package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 24000 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 45000 {
		t.Error(result)
	}
}
