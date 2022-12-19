package util

import (
	"regexp"
	"strconv"
)

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

func Group[T any](s []T, n int) (out [][]T) {
	i := 0
	group := []T{}
	for _, v := range s {
		if i != 0 && i%n == 0 {
			out = append(out, group)
			group = []T{}
		}
		group = append(group, v)
		i++
	}
	if len(group) > 0 {
		out = append(out, group)
	}
	return
}

func Split[T any](in []T, f func(T) bool) (head, tail []T) {
	appendHead := true
	for _, v := range in {
		if f(v) {
			appendHead = false
			continue
		}
		if appendHead {
			head = append(head, v)
		} else {
			tail = append(tail, v)
		}
	}
	return
}

func SlurpInt(s string) (out []int, err error) {
	for _, v := range regexp.MustCompile(`\D+`).Split(s, -1) {
		if len(v) > 0 {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			out = append(out, i)
		}
	}
	return
}

func Map[T1, T2 any](in []T1, f func(T1) T2) (out []T2) {
	for _, v := range in {
		out = append(out, f(v))
	}
	return
}

func Apply[T any](in []T, f func(T)) {
	for _, v := range in {
		f(v)
	}
}

func Filter[T any](in []T, f func(T) bool) (out []T) {
	for _, v := range in {
		if f(v) {
			out = append(out, v)
		}
	}
	return
}

func Reduce[T any](in []T, f func(a, b T) T, acc T) T {
	for _, v := range in {
		acc = f(acc, v)
	}
	return acc
}

func Min[T ~int | ~float32 | ~float64](in ...T) (out T) {
	out = in[0]
	for _, v := range in {
		if v < out {
			out = v
		}
	}
	return
}

func Max[T ~int | ~float32 | ~float64](in ...T) (out T) {
	for _, v := range in {
		if v > out {
			out = v
		}
	}
	return
}

func Cycle[T any](in []T) <-chan T {
	out := make(chan T)
	go func() {
		for {
			for _, v := range in {
				out <- v
			}
		}
	}()
	return out
}
