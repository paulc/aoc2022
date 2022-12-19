package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
	"golang.org/x/exp/slices"
)

type robot int

const (
	ore robot = iota
	clay
	obsidian
	geode
)

type cost [4]int

type state struct {
	robots    [4]int
	materials [4]int
	time      int
}

type blueprint struct {
	id    int
	costs [4]cost
}

func afford(b blueprint, s state, r robot) bool {
	return s.materials[ore] >= b.costs[r][ore] && s.materials[clay] >= b.costs[r][clay] && s.materials[obsidian] >= b.costs[r][obsidian]
}

func buy(b blueprint, s state, r robot) state {
	s.materials[ore] -= b.costs[r][ore]
	s.materials[clay] -= b.costs[r][clay]
	s.materials[obsidian] -= b.costs[r][obsidian]
	s.robots[r]++
	return s
}

func update(s state) state {
	s.materials[ore] += s.robots[ore]
	s.materials[clay] += s.robots[clay]
	s.materials[obsidian] += s.robots[obsidian]
	s.materials[geode] += s.robots[geode]
	s.time++
	return s
}

type startData struct {
	blueprints  []blueprint
	start_state state
}

func parseInput(r io.Reader) (out startData) {
	util.Must(reader.LineReader(r, func(s string) error {
		b := blueprint{}
		util.Must(fmt.Sscanf(s, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &b.id, &b.costs[ore][ore], &b.costs[clay][ore], &b.costs[obsidian][ore], &b.costs[obsidian][clay], &b.costs[geode][ore], &b.costs[geode][obsidian]))
		out.blueprints = append(out.blueprints, b)
		return nil
	}))
	out.start_state.robots[ore] = 1
	return
}

func run_blueprint(b blueprint, s state, seen set.Set[state], tmax int, out *[]int) {
	if s.time < tmax {
		if afford(b, s, geode) { // Always buy geode robots
			s2 := buy(b, update(s), geode)
			if !seen.Has(s2) {
				seen.Add(s2)
				run_blueprint(b, s2, seen, tmax, out)
			}
		} else if afford(b, s, obsidian) { // Always buy onsidian robots
			s2 := buy(b, update(s), obsidian)
			if !seen.Has(s2) {
				seen.Add(s2)
				run_blueprint(b, s2, seen, tmax, out)
			}
		} else {
			if afford(b, s, clay) {
				s2 := buy(b, update(s), clay)
				if !seen.Has(s2) {
					seen.Add(s2)
					run_blueprint(b, s2, seen, tmax, out)
				}
			}
			if afford(b, s, ore) {
				s2 := buy(b, update(s), ore)
				if !seen.Has(s2) {
					seen.Add(s2)
					run_blueprint(b, s2, seen, tmax, out)
				}
			}
			s2 := update(s)
			if !seen.Has(s2) {
				seen.Add(s2)
				run_blueprint(b, s2, seen, tmax, out)
			}
		}
	} else {
		*out = append(*out, s.materials[geode])
	}
}

func part1(input startData) (result int) {
	for _, b := range input.blueprints {
		out := []int{}
		run_blueprint(b, input.start_state, set.NewSet[state](), 24, &out)
		slices.Sort(out)
		fmt.Println(b.id, out[len(out)-1])
		result += b.id * out[len(out)-1]
	}
	return result
}

func part2(input startData) (result int) {
	result = 1
	// XXX This doesnt work on test data but does on input XXX
	for _, b := range input.blueprints[:3] {
		out := []int{}
		run_blueprint(b, input.start_state, set.NewSet[state](), 32, &out)
		slices.Sort(out)
		fmt.Println(b.id, out[len(out)-1])
		result *= out[len(out)-1]
	}
	return result
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
