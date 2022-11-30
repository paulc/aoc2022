package point

import "testing"

func TestPoint(t *testing.T) {
	p1 := Point{0, 0}
	if (p1.Move(5, 5) != Point{5, 5}) {
		t.Error(p1.Move(5, 5))
	}
	if p1.Distance(Point{5, -5}) != 10 {
		t.Error(p1.Distance(Point{5, -5}))
	}
}
