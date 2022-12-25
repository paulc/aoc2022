package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
)

type LL struct {
	val, startPos int
	prev, next    *LL
}

func (l *LL) String() string {
	out := []string{fmt.Sprintf("%d", l.val)}
	for i := l.next; i != l; i = i.next {
		out = append(out, fmt.Sprintf("%d", i.val))
	}
	return strings.Join(out, ", ")
}

func (l *LL) Move(n, max int) (out *LL) {
	out = l
	n = n % (max - 1)
	if n > 0 {
		for i := 0; i < n; i++ {
			out = out.next
		}
	} else {
		for i := 0; i <= -n; i++ {
			out = out.prev
		}
	}
	return
}

func (l *LL) MoveN(n int) (out *LL) {
	out = l
	for i := 0; i < n; i++ {
		out = out.next
	}
	return
}

func (i1 *LL) InsertAfter(i2 *LL) {
	i1.next, i1.next.prev, i2.prev, i2.next = i2, i2, i1, i1.next
}

func (i1 *LL) Remove() *LL {
	i1.prev.next, i1.next.prev, i1.next, i1.prev = i1.next, i1.prev, nil, nil
	return i1
}

type puzzle struct {
	head, zero *LL
	ptr        map[int]*LL
	length     int
}

func parseInput(r io.Reader) (out puzzle) {
	out.ptr = make(map[int]*LL)
	out.head = &LL{}
	prev, current := out.head, out.head
	i := 0
	util.Must(reader.LineReader(r, func(s string) error {
		current.val = util.Must(strconv.Atoi(s))
		current.startPos = i
		if current != prev {
			current.prev = prev
			prev.next = current
		}
		if current.val == 0 {
			out.zero = current
		}
		out.ptr[i] = current
		i++
		prev = current
		current = &LL{}
		return nil
	}))
	prev.next = out.head
	out.head.prev = prev
	out.length = i
	return
}

func part1(input puzzle) (result int) {
	for i := 0; i < input.length; i++ {
		i1 := input.ptr[i]
		i2 := i1.Move(i1.val, input.length)
		i2.InsertAfter(i1.Remove())
	}
	return input.zero.MoveN(1000).val + input.zero.MoveN(2000).val + input.zero.MoveN(3000).val
}

func part2(input puzzle) (result int) {
	for i := 0; i < input.length; i++ {
		input.ptr[i].val *= 811589153
	}
	for j := 0; j < 10; j++ {
		for i := 0; i < input.length; i++ {
			i1 := input.ptr[i]
			i2 := i1.Move(i1.val, input.length)
			i2.InsertAfter(i1.Remove())
		}
	}
	return input.zero.MoveN(1000).val + input.zero.MoveN(2000).val + input.zero.MoveN(3000).val
}

func main() {
	fmt.Println("Part1:", part1(parseInput(util.Must(os.Open("input")))))
	fmt.Println("Part2:", part2(parseInput(util.Must(os.Open("input")))))
}
