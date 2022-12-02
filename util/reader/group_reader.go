package reader

import (
	"io"
)

func GroupReader[T any](r io.Reader, groupF func(string) bool, parseF func(string) (T, error)) (out [][]T, err error) {
	i := 0
	out = append(out, []T{})
	_, err = LineReader(r, func(s string) error {
		if groupF(s) {
			out = append(out, []T{})
			i += 1
		} else {
			v, err := parseF(s)
			if err != nil {
				return err
			}
			out[i] = append(out[i], v)
		}
		return nil
	})
	return
}
