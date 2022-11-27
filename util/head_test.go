package util

import (
	"bytes"
	"strconv"
	"testing"
)

func TestHead(t *testing.T) {
	head, tail, err := Head(bytes.NewBufferString("0 0 0\n1 1 1\n2 2 2"), 1)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ArrayReader(&head, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{0, 0, 0}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}

	out, err = ArrayReader(&tail, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected = [][]int{{1, 1, 1}, {2, 2, 2}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}

func TestHeadEOF(t *testing.T) {
	head, tail, err := Head(bytes.NewBufferString("0 0 0\n1 1 1\n2 2 2"), 10)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ArrayReader(&head, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}

	if tail.Len() != 0 {
		t.Error("Tail should be nil:", tail.String())
	}
}

func TestHeadFunc(t *testing.T) {
	head, tail, err := HeadFunc(bytes.NewBufferString("0 0 0\n\n1 1 1\n2 2 2"), func(b []byte) bool { return string(b) == "\n" }, true)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ArrayReader(&head, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]int{{0, 0, 0}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}

	out, err = ArrayReader(&tail, SplitWS, strconv.Atoi)
	if err != nil {
		t.Fatal(err)
	}
	expected = [][]int{{1, 1, 1}, {2, 2, 2}}
	if !out.EqualFunc(expected, func(a, b int) bool { return a == b }) {
		t.Errorf("Out: %v\nExpected: %v\n", out, expected)
	}
}
