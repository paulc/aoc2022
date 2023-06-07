package path

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Edge[T comparable] struct {
	To   T
	Cost float64
}

type Graph[T comparable] map[T][]Edge[T]

func (g Graph[T]) String() string {
	out := []string{}
	for k, v := range g {
		out = append(out, fmt.Sprintf("%v -> %v", k, v))
	}
	slices.Sort(out)
	return strings.Join(out, "\n")
}
