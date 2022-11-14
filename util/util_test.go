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
	if slices.CompareFunc(out, expected, func(e1, e2 []int) int { return slices.Compare(e1, e2) }) != 0 {
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
	if slices.CompareFunc(out, expected, func(e1, e2 []int) int { return slices.Compare(e1, e2) }) != 0 {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestArrayReaderPoint(t *testing.T) {

	r := bytes.NewBufferString(data_point)
	out, err := ArrayReader(r, MakeStringSplitter(","), ParsePoint)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
}
