package main

import (
	"fmt"

	"github.com/paulc/aoc2022/util/array"
)

type edge struct {
	next_face, next_edge int
}

type face struct {
	tiles          array.Array[tile]
	startx, starty int
	w, h           int
	edges          [4]edge
}

type cube [6]face

type cubepos struct {
	f, x, y int
	dir     facing
}

func (p cubepos) String() string {
	return fmt.Sprintf("Face: %d (%d,%d) %s", p.f, p.x, p.y, p.dir)
}

func extractFace(in array.Array[tile], ix, iy, w, h int) (out face) {
	out.tiles = make(array.Array[tile], h)
	out.w = w
	out.h = h
	out.startx, out.starty = ix*w, iy*h
	for y := 0; y < w; y++ {
		out.tiles[y] = in[iy*h+y][ix*w : ix*w+w]
	}
	return
}

func mapCube(in array.Array[tile], m face_map) (out cube) {
	h := len(in) / m.ny
	w := len(in[0]) / m.nx
	// Extract faces
	for i, v := range m.faces {
		out[i] = extractFace(in, v[0], v[1], w, h)
		out[i].edges = m.edges[i]
	}
	return
}

func printCube(c cube, p cubepos) {
	for i := 0; i < 6; i++ {
		if p.f == i {
			fmt.Println(p)
			t := c[i].tiles.Copy()
			t.Set(p.x, p.y, tile(int(right)+int(p.dir)))
			fmt.Println(t)
		}
	}
}

func (c cube) Move(p cubepos, m move) cubepos {
	if m.turn {
		if m.direction == "R" {
			p.dir = (p.dir + 1) % nFacing
		} else {
			p.dir = (p.dir + 3) % nFacing
		}
	} else {
		next := p
		w := c[p.f].w - 1
		for i := 0; i < m.count; i++ {
			delta := dxy[p.dir]
			next.x, next.y = p.x+delta.dx, p.y+delta.dy
			if next.x < 0 {
				e := c[next.f].edges[3]
				next.f = e.next_face
				switch e.next_edge {
				case 0:
					next.x, next.y, next.dir = next.y, 0, D
				case 1:
					next.x, next.y, next.dir = w, next.y, L
				case 2:
					next.x, next.y, next.dir = w-next.y, w, U
				case 3:
					next.x, next.y, next.dir = 0, w-next.y, R
				}
			} else if next.x > w {
				e := c[next.f].edges[1]
				next.f = e.next_face
				switch e.next_edge {
				case 0:
					next.x, next.y, next.dir = w-next.y, 0, D
				case 1:
					next.x, next.y, next.dir = w, w-next.y, L
				case 2:
					next.x, next.y, next.dir = next.y, w, U
				case 3:
					next.x, next.y, next.dir = 0, next.y, R
				}
			} else if next.y < 0 {
				e := c[next.f].edges[0]
				next.f = e.next_face
				switch e.next_edge {
				case 0:
					next.x, next.y, next.dir = next.x, 0, D
				case 1:
					next.x, next.y, next.dir = w, w-next.x, L
				case 2:
					next.x, next.y, next.dir = next.x, w, U
				case 3:
					next.x, next.y, next.dir = 0, next.x, R
				}
			} else if next.y > w {
				e := c[next.f].edges[2]
				next.f = e.next_face
				switch e.next_edge {
				case 0:
					next.x, next.y, next.dir = next.x, 0, D
				case 1:
					next.x, next.y, next.dir = w, next.x, L
				case 2:
					next.x, next.y, next.dir = w-next.x, w, U
				case 3:
					next.x, next.y, next.dir = 0, w-next.x, R
				}
			}
			switch c[next.f].tiles[next.y][next.x] {
			case solid:
				return p // Blocked
			case open:
				p = next
			}
		}
	}
	return p
}

type face_map struct {
	nx, ny int
	start  cubepos
	faces  [][2]int
	edges  [6][4]edge
}

var part2_map = face_map{
	nx:    3,
	ny:    4,
	start: cubepos{0, 0, 0, R},
	faces: [][2]int{{1, 0}, {2, 0}, {1, 1}, {0, 2}, {1, 2}, {0, 3}},
	edges: [6][4]edge{
		{{5, 3}, {1, 3}, {2, 0}, {3, 3}},
		{{5, 2}, {4, 1}, {2, 1}, {0, 1}},
		{{0, 2}, {1, 2}, {4, 0}, {3, 0}},
		{{2, 3}, {4, 3}, {5, 0}, {0, 3}},
		{{2, 2}, {1, 1}, {5, 1}, {3, 1}},
		{{3, 2}, {4, 2}, {1, 0}, {0, 0}},
	},
}

func part2(input puzzle, m face_map) (result int) {
	cube := mapCube(input.cave, m)
	pos := m.start
	for _, move := range input.moves {
		pos = cube.Move(pos, move)
	}
	end_face := cube[pos.f]
	result = 4*(end_face.startx+pos.x+1) + 1000*(end_face.starty+pos.y+1) + int(pos.dir)
	return
}
