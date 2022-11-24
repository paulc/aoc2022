package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

func Head(r io.Reader, n int) (head bytes.Buffer, tail bytes.Buffer, err error) {
	b := bufio.NewReader(r)
	for i := 0; i < n; i++ {
		var line []byte
		line, err := b.ReadBytes('\n')
		head.Write(line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return head, tail, nil
			} else {
				return head, tail, err
			}
		}
	}
	t, err := io.ReadAll(b)
	if err != nil {
		return head, tail, err
	}
	tail.Write(t)
	return head, tail, nil
}

func HeadFunc(r io.Reader, splitf func([]byte) bool, skipMatch bool) (head bytes.Buffer, tail bytes.Buffer, err error) {
	b := bufio.NewReader(r)
	for {
		var line []byte
		line, err := b.ReadBytes('\n')
		if splitf(line) {
			if !skipMatch {
				tail.Write(line)
			}
			break
		}
		head.Write(line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return head, tail, nil
			} else {
				return head, tail, err
			}
		}
	}
	t, err := io.ReadAll(b)
	if err != nil {
		return head, tail, err
	}
	tail.Write(t)
	return head, tail, nil
}
