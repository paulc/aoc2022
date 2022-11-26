package util

import (
	"errors"
	"fmt"
	"strings"
)

type Point struct {
	X, Y int
}

func (p Point) Move(dx, dy int) Point {
	return Point{p.X + dx, p.Y + dy}
}

func absint(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func (p Point) Distance(p2 Point) int {
	return absint(p.X-p2.X) + absint(p.Y-p2.Y)
}

type Grid[T any] struct {
	X0, Y0, X1, Y1 int
	Width, Height  int
	Data           []T
}

func NewGrid[T any](x0, y0, x1, y1 int) (*Grid[T], error) {
	if x1 <= x0 || y1 <= y0 {
		return nil, errors.New("Invalid bounds")
	}
	g := &Grid[T]{X0: x0, Y0: y0, X1: x1, Y1: y1}
	g.Width = x1 - x0 + 1
	g.Height = y1 - y0 + 1
	g.Data = make([]T, g.Width*g.Height)
	return g, nil
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

func (g *Grid[T]) AdjacentWrap(p Point) (out []Point) {
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
		out = append(out, g.Move(p, v.dx, v.dy))
	}
	return
}

func (g *Grid[T]) Move(p Point, dx, dy int) Point {
	p1 := p.Move(dx, dy)
	if !g.CheckBounds(p1) {
		if p1.X > g.X1 {
			p1.X = g.X0 + (p1.X-g.X1-1)%g.Width
		} else if p1.X < g.X0 {
			p1.X = g.X1 - (g.X0-p1.X-1)%g.Width
		}
		if p1.Y > g.Y1 {
			p1.Y = g.Y0 + (p1.Y-g.Y1-1)%g.Width
		} else if p1.Y < g.Y0 {
			p1.Y = g.Y1 - (g.Y0-p1.Y-1)%g.Width
		}
	}
	return p1
}
