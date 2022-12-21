package main

import (
	"bytes"
	"strings"
	"testing"
)

const testdata = `
root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32
`

var input = parseInput(bytes.NewBufferString(strings.TrimSpace(testdata)))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 152 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 301 {
		t.Error(result)
	}
}
