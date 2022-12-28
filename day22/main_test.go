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

var test_map = face_map{
	nx:    4,
	ny:    3,
	start: cubepos{0, 0, 0, R},
	faces: [][2]int{{2, 0}, {0, 1}, {1, 1}, {2, 1}, {2, 2}, {3, 2}},
	edges: [6][4]edge{
		{{1, 0}, {5, 1}, {3, 0}, {2, 0}},
		{{0, 0}, {2, 3}, {4, 2}, {5, 2}},
		{{0, 3}, {3, 3}, {4, 3}, {1, 1}},
		{{0, 2}, {5, 0}, {4, 0}, {2, 1}},
		{{3, 2}, {5, 3}, {1, 2}, {2, 2}},
		{{3, 1}, {0, 1}, {1, 3}, {4, 1}},
	},
}

func TestPart2(t *testing.T) {
	result := part2(input, test_map)
	if result != 5031 {
		t.Error(result)
	}
}
