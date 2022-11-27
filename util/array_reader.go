package util

import (
	"io"
	"regexp"
	"strings"
)

func ArrayReaderFunc[T any](r io.Reader, f func(s string) ([]T, error)) (out Array[T], err error) {
	_, err = LineReader(r, func(s string) error {
		if len(s) > 0 {
			v, err := f(s)
			if err != nil {
				return err
			}
			out = append(out, v)
		}
		return nil
	})
	return
}

func LineParser[T any](line string, splitF func(string) ([]string, error), parseF func(string) (T, error)) (out []T, err error) {
	if len(line) > 0 {
		var split []string
		split, err = splitF(line)
		if err != nil {
			return
		}
		for _, v := range split {
			var p T
			p, err = parseF(v)
			if err != nil {
				return
			}
			out = append(out, p)
		}
	}
	return
}

func ArrayReader[T any](r io.Reader, splitF func(string) ([]string, error), parseF func(string) (T, error)) (out Array[T], err error) {
	_, err = LineReader(r, func(s string) error {
		if len(s) > 0 {
			line, err := LineParser(s, splitF, parseF)
			if err != nil {
				return err
			}
			out = append(out, line)
		}
		return nil
	})
	return
}

var ws = regexp.MustCompile(`\s+`)

func SplitWS(s string) ([]string, error) {
	return ws.Split(s, -1), nil
}

func MakeStringSplitter(sep string) func(s string) ([]string, error) {
	return func(s string) ([]string, error) {
		return strings.Split(s, sep), nil
	}
}

func MakeRegexpSplitter(re string) func(s string) ([]string, error) {
	re_c := regexp.MustCompile(re)
	return func(s string) ([]string, error) {
		return re_c.Split(s, -1), nil
	}
}
