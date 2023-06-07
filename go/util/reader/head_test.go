package reader

import (
	"bytes"
	"testing"
)

func TestHead(t *testing.T) {
	head, tail, err := Head(bytes.NewBufferString("0 0 0\n1 1 1\n2 2 2"), 1)
	if err != nil {
		t.Fatal(err)
	}
	if head.String() != "0 0 0\n" {
		t.Error(head.String())
	}
	if tail.String() != "1 1 1\n2 2 2" {
		t.Error(tail.String())
	}
}

func TestHeadEOF(t *testing.T) {
	head, tail, err := Head(bytes.NewBufferString("0 0 0\n1 1 1\n2 2 2"), 10)
	if err != nil {
		t.Fatal(err)
	}
	if head.String() != "0 0 0\n1 1 1\n2 2 2" {
		t.Error(head.String())
	}
	if tail.String() != "" {
		t.Error(tail.String())
	}
}

func TestHeadFunc(t *testing.T) {
	head, tail, err := HeadFunc(bytes.NewBufferString("0 0 0\n\n1 1 1\n2 2 2"), func(b []byte) bool { return string(b) == "\n" }, true)
	if err != nil {
		t.Fatal(err)
	}
	if head.String() != "0 0 0\n" {
		t.Error(head.String())
	}
	if tail.String() != "1 1 1\n2 2 2" {
		t.Error(tail.String())
	}
}
