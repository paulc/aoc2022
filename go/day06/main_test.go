package main

import (
	"bytes"
	"testing"
)

var testdata = []struct {
	b      []byte
	p1, p2 int
}{
	{bytes.NewBufferString("mjqjpqmgbljsphdztnvjfqwrcgsmlb").Bytes(), 7, 19},
	{bytes.NewBufferString("bvwbjplbgvbhsrlpgdmjqwftvncz").Bytes(), 5, 23},
	{bytes.NewBufferString("nppdvjthqldpwncqszvftbrmjlhg").Bytes(), 6, 23},
	{bytes.NewBufferString("nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg").Bytes(), 10, 29},
	{bytes.NewBufferString("zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw").Bytes(), 11, 26},
}

func TestPart1(t *testing.T) {
	for _, v := range testdata {
		if part1(v.b) != v.p1 {
			t.Error(string(v.b), v.p1, part1(v.b))
		}
	}
}

func TestPart2(t *testing.T) {
	for _, v := range testdata {
		if part2(v.b) != v.p2 {
			t.Error(string(v.b), v.p2, part2(v.b))
		}
	}
}
