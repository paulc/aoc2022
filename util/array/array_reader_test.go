package array

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

type _point struct {
	X, Y int
}

func ParsePoint(s string) (p _point, err error) {
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

func TestLineParser(t *testing.T) {
	out, err := LineParser("1 2 3", SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(out, []int{1, 2, 3}) {
		t.Fatal(err)
	}
}

func TestLineParserErr(t *testing.T) {
	_, err := LineParser("1 xx 3", SplitWS, strconv.Atoi)
	if err == nil {
		t.Fatal("Expected parse error")
	}
}

func TestArrayReader(t *testing.T) {
	r := bytes.NewBufferString(data)
	out, err := ArrayReader(r, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
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
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestArrayReaderPoint(t *testing.T) {
	r := bytes.NewBufferString(data_point)
	out, err := ArrayReader(r, MakeStringSplitter(","), ParsePoint)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]_point{{_point{1, 2}, _point{3, 4}, _point{5, 6}}, {_point{2, 1}, _point{4, 3}, _point{6, 5}}}
	if !out.EqualFunc(expected, func(a, b _point) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestArrayReaderPointErr(t *testing.T) {
	r := bytes.NewBufferString("1:2,xx,3:4")
	_, err := ArrayReader(r, MakeStringSplitter(","), ParsePoint)
	if err == nil {
		t.Fatal("Expected parse error")
	}
}
