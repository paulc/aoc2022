package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

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

func best_estimate(input cave, tnow, tmax int, location string, available set.Set[string]) (result int) {
	rates := []int{}
	fastest := tmax
	available.Apply(func(s string) {
		rates = append(rates, input.valveMap[s])
		if c := input.costs[location+s]; c < fastest {
			fastest = c
		}
	})
	slices.Sort(rates)
	for tleft := tmax - tnow - fastest - 1; tleft > 0; tleft -= 2 {
		if len(rates) > 0 {
			result += rates[len(rates)-1] * tleft
			rates = rates[:len(rates)-1]
		}
	}
	return
}

func search(input cave, current state, available set.Set[string], visited set.Set[state], tmax int, best *int) {
	visited.Add(current)
	if current.pressure > *best {
		*best = current.pressure
	}
	if current.pressure+best_estimate(input, current.time, tmax, current.location, available) > *best {
		for _, v := range available.Keys() {
			t := current.time + input.costs[current.location+v] + 1
			if t < tmax {
				next := state{t, v, addValve(current.valvesOn, v), current.pressure + (input.valveMap[v] * (tmax - t))}
				if !visited.Has(next) {
					search(input, next, available.Copy().Remove(v), visited, tmax, best)
				}
			}
		}
	}
}

func part1(input cave) (result int) {

	visited := set.NewSet[state]()
	start := state{0, "AA", "", 0}
	visited.Add(start)

	best := 0
	search(input, start, input.available, visited, 30, &best)
	return best
}

func calculatePair(input cave, split []string) int {
	start := state{0, "AA", "", 0}
	e := set.NewSetFrom(split)
	p := input.available.Difference(e)
	best_e := 0
	best_p := 0
	visited_e := set.NewSetFrom([]state{start})
	visited_p := set.NewSetFrom([]state{start})

	search(input, start, e, visited_e, 26, &best_e)
	search(input, start, p, visited_p, 26, &best_p)

	return best_e + best_p
}

func part2(input cave) (result int) {

	split := input.available.Len() / 2
	keys := input.available.Keys()
	results := make(chan int)
	done := make(chan bool)
	wg := sync.WaitGroup{}

	for i := 1; i <= split; i++ {
		for _, v := range util.Combinations(keys, i) {
			wg.Add(1)
			go func(s []string) {
				defer wg.Done()
				results <- calculatePair(input, s)
			}(v)
		}
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	running := true
	for running {
		select {
		case v := <-results:
			if v > result {
				result = v
			}
		case <-done:
			running = false
		}
	}
	return
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
