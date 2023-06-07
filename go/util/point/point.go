package point

type Point struct {
	X, Y int
}

func (p Point) Move(dx, dy int) Point {
	return Point{p.X + dx, p.Y + dy}
}

func (p Point) Distance(p2 Point) int {
	return absint(p.X-p2.X) + absint(p.Y-p2.Y)
}

func (p Point) Xdistance(p2 Point) int {
	return absint(p.X - p2.X)
}

func (p Point) Ydistance(p2 Point) int {
	return absint(p.Y - p2.Y)
}

func (p Point) Adjacent() (out []Point) {
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}} {
		out = append(out, Point{p.X + v.dx, p.Y + v.dy})
	}
	return
}

func (p Point) AdjacentDiagonal() (out []Point) {
	for _, v := range []struct{ dx, dy int }{{-1, 0}, {0, -1}, {1, 0}, {0, 1}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}} {
		out = append(out, Point{p.X + v.dx, p.Y + v.dy})
	}
	return
}

func absint(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
