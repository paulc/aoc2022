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

type stateElephant struct {
	timeYou          int
	timeElephant     int
	locationYou      string
	locationElephant string
	valvesOn         string
	pressure         int
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

func pairs[T any](in []T) (out [][2]T) {
	for i := range in {
		for j := range in {
			if j != i {
				out = append(out, [2]T{in[i], in[j]})
			}
		}
	}
	return
}

func searchElephant(input cave, current stateElephant, available set.Set[string], visited set.Set[stateElephant], tmax int) {
	visited.Add(current)
	if available.Len() == 1 {
		v := available.Keys()[0]
		next := current
		tyou := current.timeYou + input.costs[current.locationYou+v] + 1
		telephant := current.timeElephant + input.costs[current.locationElephant+v] + 1
		remaining := available.Copy()
		if tyou < tmax && tyou < telephant {
			next.timeYou = tyou
			next.locationYou = v
			next.pressure += input.valveMap[v] * (tmax - tyou)
			next.valvesOn = addValve(next.valvesOn, v)
			remaining.Remove(v)
		} else if telephant < tmax {
			next.timeElephant = telephant
			next.locationElephant = v
			next.pressure += input.valveMap[v] * (tmax - telephant)
			next.valvesOn = addValve(next.valvesOn, v)
			remaining.Remove(v)
		}
		if !visited.Has(next) {
			searchElephant(input, next, remaining, visited, tmax)
		}

		/*
			v := available.Keys()[0]
			tyou := current.timeYou + input.costs[current.locationYou+v] + 1
			lyou := current.locationYou
			telephant := current.timeElephant + input.costs[current.locationElephant+v] + 1
			lelephant := current.locationElephant
			pressure := current.pressure
			von := current.valvesOn
			if tyou <= telephant && tyou < tmax {
				lyou = v
				pressure += (input.valveMap[lyou] * (tmax - tyou))
				von = addValve(von, lyou)
				next := stateElephant{tyou, telephant, lyou, lelephant, von, pressure}
				visited.Add(next)
			} else if telephant < tyou && telephant < tmax {
				lelephant = v
				pressure += (input.valveMap[lelephant] * (tmax - telephant))
				von = addValve(von, lelephant)
				next := stateElephant{tyou, telephant, lyou, lelephant, von, pressure}
				visited.Add(next)
			}
		*/
	} else {
		for _, v := range pairs(available.Keys()) {
			next := current
			tyou := current.timeYou + input.costs[current.locationYou+v[0]] + 1
			telephant := current.timeElephant + input.costs[current.locationElephant+v[1]] + 1
			remaining := available.Copy()
			if tyou < tmax {
				next.timeYou = tyou
				next.locationYou = v[0]
				next.pressure += input.valveMap[v[0]] * (tmax - tyou)
				next.valvesOn = addValve(next.valvesOn, v[0])
				remaining.Remove(v[0])
			}
			if telephant < tmax {
				next.timeElephant = telephant
				next.locationElephant = v[1]
				next.pressure += input.valveMap[v[1]] * (tmax - telephant)
				next.valvesOn = addValve(next.valvesOn, v[1])
				remaining.Remove(v[1])
			}
			if !visited.Has(next) {
				searchElephant(input, next, remaining, visited, tmax)
			}
		}
	}
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
	visited := set.NewSet[stateElephant]()
	start := stateElephant{0, 0, "AA", "AA", "", 0}
	visited.Add(start)

	searchElephant(input, start, input.available, visited, 26)

	out := visited.Keys()
	slices.SortFunc(out, func(a, b stateElephant) bool { return a.pressure < b.pressure })

	return out[len(out)-1].pressure
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
