package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
1
2
-3
3
-2
0
4
`

func TestPart1(t *testing.T) {
	result := part1(parseInput(bytes.NewBufferString(strings.TrimSpace(testdata))))
	if result != 3 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(parseInput(bytes.NewBufferString(strings.TrimSpace(testdata))))
	if result != 0 {
		t.Error(result)
	}
}
