package main

import (
	"github.com/paulc/aoc2022/util/set"
	"golang.org/x/exp/slices"
)

type stateElephant struct {
	timeYou          int
	timeElephant     int
	locationYou      string
	locationElephant string
	valvesOn         string
	pressure         int
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

func part2_elephant(input cave) (result int) {
	visited := set.NewSet[stateElephant]()
	start := stateElephant{0, 0, "AA", "AA", "", 0}
	visited.Add(start)

	searchElephant(input, start, input.available, visited, 26)

	out := visited.Keys()
	slices.SortFunc(out, func(a, b stateElephant) bool { return a.pressure < b.pressure })

	return out[len(out)-1].pressure
}
