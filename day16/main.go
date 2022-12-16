package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/path"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
	"golang.org/x/exp/slices"
)

type cave struct {
	valveMap map[string]int
	costs    map[string]int
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
	util.Apply(interesting, func(v valve) {
		routes = append(routes, v.key)
		valveMap[v.key] = v.flow
	})
	for _, start := range routes {
		for _, r := range paths.AstarMultiple(start, routes, func(s string) float64 { return 1.0 }) {
			costs[start+r.End] = int(r.Cost)
		}
	}
	return cave{valveMap, costs}
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
		if t < tmax && available.Len() > 0 {
			next := state{t, v, addValve(current.valvesOn, v), current.pressure + (input.valveMap[v] * (tmax - t))}
			if !visited.Has(next) {
				search(input, next, available.Copy().Remove(v), visited, tmax)
			}
		}
	}
}

func part1(input cave) (result int) {
	available := set.NewSet[string]()
	for k, _ := range input.valveMap {
		available.Add(k)
	}

	visited := set.NewSet[state]()
	start := state{0, "AA", "", 0}
	visited.Add(start)

	search(input, start, available, visited, 30)
	out := visited.Keys()
	slices.SortFunc(out, func(a, b state) bool { return a.pressure < b.pressure })
	return out[len(out)-1].pressure
}

func part2(input cave) (result int) {
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
