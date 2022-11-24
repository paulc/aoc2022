package util

import (
	"bytes"
	"strconv"
	"testing"
)

var testReader = bytes.NewBufferString(`0 0 0
1 1 1
2 2 2`)

func TestHead(t *testing.T) {
	head, tail, err := Head(testReader, 1)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ArrayReader(&head, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{0, 0, 0}}
	if !ArrayEquals(out, expected) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}

	out, err = ArrayReader(&tail, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected = [][]int{{1, 1, 1}, {2, 2, 2}}
	if !ArrayEquals(out, expected) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}

}
