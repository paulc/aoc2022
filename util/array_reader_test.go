package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

const data = `
1 2 3
4 5 6
7 8 9
`

const data_point = `
1:2,3:4,5:6
2:1,4:3,6:5
`

type Point struct {
	X, Y int
}

func ArrayEquals[T comparable](a1, a2 [][]T) bool {
	return slices.EqualFunc(a1, a2, func(s1, s2 []T) bool { return slices.Equal(s1, s2) })
}

func ParsePoint(s string) (p Point, err error) {
	xy := strings.SplitN(s, ":", 2)
	if len(xy) != 2 {
		err = fmt.Errorf("Invalid Point: %s", s)
		return
	}
	p.X, err = strconv.Atoi(xy[0])
	if err != nil {
		return
	}
	p.Y, err = strconv.Atoi(xy[1])
	return
}

func TestArrayReader(t *testing.T) {

	r := bytes.NewBufferString(data)
	out, err := ArrayReader(r, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	if !ArrayEquals(out, expected) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestArrayReaderFunc(t *testing.T) {

	r := bytes.NewBufferString(data)
	out, err := ArrayReaderFunc(r, func(s string) (out []int, err error) {
		strings.Split(s, " ")
		for _, v := range strings.Split(s, " ") {
			var i int
			i, err = strconv.Atoi(v)
			if err != nil {
				return
			}
			out = append(out, i)
		}
		return
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	if !ArrayEquals(out, expected) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestArrayReaderPoint(t *testing.T) {

	r := bytes.NewBufferString(data_point)
	out, err := ArrayReader(r, MakeStringSplitter(","), ParsePoint)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]Point{{Point{1, 2}, Point{3, 4}, Point{5, 6}}, {Point{2, 1}, Point{4, 3}, Point{6, 5}}}
	if !ArrayEquals(out, expected) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}
