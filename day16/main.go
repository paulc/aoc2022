package main

import (
	"container/heap"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/path"
	"github.com/paulc/aoc2022/util/priqueue"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
	"golang.org/x/exp/slices"
)

type cave struct {
	valveMap  map[string]int
	costs     map[string]int
	available set.Set[string]
}

type valve struct {
	key   string
	flow  int
	paths []string
}

type state struct {
	time     int
	location string
	valvesOn string
	pressure int
}

func parseInput(r io.Reader) cave {
	paths := make(path.Graph[string])
	valves := util.Map(util.Must(reader.Lines(r)), func(s string) (out valve) {
		util.Must(fmt.Sscanf(s, "Valve %s has flow rate=%d;", &out.key, &out.flow))
		out.paths = strings.Split(strings.SplitN(s, " ", 10)[9], ", ")
		paths[out.key] = util.Map(out.paths, func(e string) path.Edge[string] {
			return path.Edge[string]{e, 1.0}
		})
		return
	})
	interesting := util.Filter(valves, func(v valve) bool { return v.flow > 0 })
	routes := []string{"AA"}
	costs := make(map[string]int)
	valveMap := make(map[string]int)
	available := set.NewSet[string]()
	util.Apply(interesting, func(v valve) {
		routes = append(routes, v.key)
		valveMap[v.key] = v.flow
		available.Add(v.key)
	})
	for _, start := range routes {
		for _, r := range paths.AstarMultiple(start, routes, func(s string) float64 { return 1.0 }) {
			costs[start+r.End] = int(r.Cost)
		}
	}
	return cave{valveMap, costs, available}
}

func addValve(on, v string) string {
	s := set.NewSetFrom(strings.Split(on, ","))
	s.Add(v)
	keys := s.Keys()
	slices.Sort(keys)
	return strings.Join(keys, ",")
}

func search(input cave, current state, available set.Set[string], visited set.Set[state], tmax int) {
	visited.Add(current)
	for _, v := range available.Keys() {
		t := current.time + input.costs[current.location+v] + 1
		if t < tmax {
			next := state{t, v, addValve(current.valvesOn, v), current.pressure + (input.valveMap[v] * (tmax - t))}
			if !visited.Has(next) {
				search(input, next, available.Copy().Remove(v), visited, tmax)
			}
		}
	}
}

func split[T comparable](in set.Set[T], n int) (out []struct{ a, b set.Set[T] }) {
	if n == 1 {
		keys := in.Keys()
		for i := 0; i < len(keys); i++ {
			a := set.NewSetFrom([]T{keys[i]})
			b := in.Copy().Remove(keys[i])
			out = append(out, struct{ a, b set.Set[T] }{a, b})
		}
		return
	}
	for _, v := range split(in, n-1) {
		for _, v2 := range split(v.b, 1) {
			a := v.a.Union(v2.a)
			b := v2.b
			out = append(out, struct{ a, b set.Set[T] }{a, b})
		}
	}
	return
}

func search2(input cave, start state, startValves set.Set[string], tmax int) (result int) {
	q := priqueue.NewPriorityQueue[state](func(s state) float64 { return float64(-s.pressure) })
	visited := set.NewSetFrom([]state{start})
	heap.Push(q, start)
	for q.Len() > 0 {
		current := heap.Pop(q).(state)
		if current.pressure > result {
			result = current.pressure
		}
		available := startValves.Difference(set.NewSetFrom(strings.Split(current.valvesOn, ",")))
		for v := range available {
			tnext := current.time + input.costs[current.location+v] + 1
			if tnext < tmax {
				next := state{tnext, v, addValve(current.valvesOn, v), current.pressure + (input.valveMap[v] * (tmax - tnext))}
				if !visited.Has(next) {
					visited.Add(next)
					heap.Push(q, next)
				}
			}
		}
	}
	return
}

func part1(input cave) (result int) {

	visited := set.NewSet[state]()
	start := state{0, "AA", "", 0}
	visited.Add(start)

	search(input, start, input.available, visited, 30)

	out := visited.Keys()
	slices.SortFunc(out, func(a, b state) bool { return a.pressure < b.pressure })

	return out[len(out)-1].pressure
}

func part2(input cave) (result int) {
	//start := state{0, "AA", "", 0}
	// mid := len(input.available) / 2
	for _, v := range split(input.available, 5) {
		fmt.Println(v.a.Keys(), v.b.Keys())
		/*
			r := search2(input, start, v.a, 26) + search2(input, start, v.b, 26)
			fmt.Println(v.a.Keys(), v.b.Keys(), r)
			if r > result {
				result = r
				fmt.Println(v.a.Keys(), v.b.Keys(), result)
			}
		*/
	}
	return
}

/*

	// XXX This doesnt work properly (but does get the right answer) XXX

	visited := set.NewSet[state]()
	start := state{0, "AA", "", 0}
	visited.Add(start)

	elephant_avail := input.available.Copy()
	elephant_visited := visited.Copy()
	search(input, start, input.available, visited, 26)

	out := visited.Keys()
	slices.SortFunc(out, func(a, b state) bool { return a.pressure < b.pressure })
	result = out[len(out)-1].pressure

	done := set.NewSetFrom(strings.Split(out[len(out)-1].valvesOn, ","))

	for k := range done {
		elephant_avail.Remove(k)
	}

	search(input, start, elephant_avail, elephant_visited, 26)

	elephant_out := elephant_visited.Keys()
	slices.SortFunc(elephant_out, func(a, b state) bool { return a.pressure < b.pressure })
	result += elephant_out[len(elephant_out)-1].pressure
*/

func main() {
	input := parseInput(util.Must(os.Open("input")))
	// fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
