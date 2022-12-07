package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type Dir struct {
	parent   *Dir
	name     string
	children map[string]*Dir
	files    map[string]int
	size     int
}

func NewDir(parent *Dir, name string) *Dir {
	return &Dir{parent: parent, name: name, children: make(map[string]*Dir), files: make(map[string]int), size: -1}
}

func (d *Dir) Walk(f func(d *Dir)) {
	for _, v := range d.children {
		v.Walk(f)
	}
	f(d)
}

func (d *Dir) Size() (size int) {
	for _, v := range d.files {
		size += v
	}
	for _, v := range d.children {
		size += v.Size()
	}
	return
}

func (d *Dir) Cd(name string) *Dir {
	_, found := d.children[name]
	if !found {
		d.children[name] = NewDir(d, name)
	}
	return d.children[name]
}

func parseInput(r io.Reader) *Dir {
	root := NewDir(nil, "/")
	cwd := root
	for _, v := range util.Must(reader.Lines(r)) {
		line := strings.Split(v, " ")
		if line[0] == "$" {
			if line[1] == "cd" {
				switch line[2] {
				case "/":
					cwd = root
				case "..":
					cwd = cwd.parent
				default:
					cwd = cwd.Cd(line[2])
				}
			}
		} else { // Dir list
			if line[0] == "dir" {
				cwd.Cd(line[1])
			} else {
				cwd.files[line[1]] = util.Must(strconv.Atoi(line[0]))
			}
		}
	}
	return root
}

func part1(root *Dir) (result int) {
	root.Walk(func(d *Dir) {
		if s := d.Size(); s < 100000 {
			result += s
		}
	})
	return result
}

func part2(root *Dir) (result int) {
	need := 30000000 - (70000000 - root.Size())
	avail := []int{}
	root.Walk(func(d *Dir) {
		if s := d.Size(); s > need {
			avail = append(avail, s)
		}
	})
	sort.Ints(avail)
	return avail[0]
}

func main() {
	root := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(root))
	fmt.Println("Part2:", part2(root))
}
