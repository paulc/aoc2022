package util

import (
	"fmt"
	"strings"
)

type Point struct {
	X, Y int
}

func (p Point) Move(dx, dy int) Point {
	return Point{p.X + dx, p.Y + dy}
}

type Grid[T any] struct {
	X0, Y0, X1, Y1 int
	Width, Height  int
	Data           []T
}

func NewGrid[T any](x0, y0, x1, y1 int) *Grid[T] {
	g := &Grid[T]{}
	g.X0 = x0
	g.Y0 = y0
	g.X1 = x1
	g.Y1 = y1
	g.Width = x1 - x0
	g.Height = y1 - y0
	g.Data = make([]T, g.Width*g.Height)
	return g
}

func (g *Grid[T]) CheckBounds(p Point) bool {
	return !(p.X < g.X0 || p.X > g.X1 || p.Y < g.Y0 || p.Y > g.Y1)
}

func (g *Grid[T]) Set(p Point, val T) {
	// We sliently ignore out of bounds errors
	if !g.CheckBounds(p) {
		return
	}
	g.Data[(p.X-g.X0)+(p.Y-g.Y0)*g.Width] = val
}

func (g *Grid[T]) Get(p Point) (out T) {
	// Return zero val if out of bounds
	if !g.CheckBounds(p) {
		return
	}
	return g.Data[(p.X-g.X0)+(p.Y-g.Y0)*g.Width]
}

func (g *Grid[T]) String() string {
	rows := make([]string, g.Height)
	for y := 0; y < g.Height; y++ {
		line := make([]string, g.Width)
		for x := 0; x < g.Width; x++ {
			line[x] = fmt.Sprintf("%v", g.Data[(y*g.Width)+x])
		}
		rows[y] = strings.Join(line, " ")
	}
	return strings.Join(rows, "\n")
}

func (g *Grid[T]) Adjacent(p Point) (out []Point) {
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
		p1 := p.Move(v.dx, v.dy)
		if g.CheckBounds(p1) {
			out = append(out, p1)
		}
	}
	return
}
