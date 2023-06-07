package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 0 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 0 {
		t.Error(result)
	}
}
