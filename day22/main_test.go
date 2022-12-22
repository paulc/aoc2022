package main

import (
	"bytes"
	"testing"
)

const testdata = `        ...#    
        .#..    
        #...    
        ....    
...#.......#    
........#...    
..#....#....    
..........#.    
        ...#....
        .....#..
        .#......
        ......#.

10R5L5R10L4R5L5`

var input = parseInput(bytes.NewBufferString(testdata))

func TestPart1(t *testing.T) {
	result := part1(input)
	if result != 6032 {
		t.Error(result)
	}
}

func TestPart2(t *testing.T) {
	result := part2(input)
	if result != 0 {
		t.Error(result)
	}
}
