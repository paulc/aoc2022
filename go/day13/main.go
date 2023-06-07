package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/reader"
	"golang.org/x/exp/slices"
)

func parseInput(r io.Reader) (out []interface{}) {
	for _, v := range util.Must(reader.Lines(r)) {
		if v != "" {
			var p interface{}
			json.Unmarshal([]byte(v), &p)
			out = append(out, p)
		}
	}
	return
}

func checkOrder(l, r interface{}) int {
	switch lt := l.(type) {
	case float64:
		switch rt := r.(type) {
		case float64:
			return int(rt) - int(lt)
		case []interface{}:
			return checkOrder([]interface{}{lt}, rt)
		}
	case []interface{}:
		switch rt := r.(type) {
		case float64:
			return checkOrder(lt, []interface{}{rt})
		case []interface{}:
			for i := 0; i < util.Min(len(lt), len(rt)); i++ {
				if o := checkOrder(lt[i], rt[i]); o != 0 {
					return o
				}
			}
			return len(rt) - len(lt)
		}
	}
	return 0
}

func part1(input []interface{}) (result int) {
	for i := 0; i < len(input); i += 2 {
		if checkOrder(input[i], input[i+1]) > 0 {
			result += (i / 2) + 1
		}
	}
	return result
}

func part2(input []interface{}) (result int) {
	var d1, d2 interface{}
	json.Unmarshal([]byte("[[2]]"), &d1)
	json.Unmarshal([]byte("[[6]]"), &d2)
	input = append(input, d1, d2)
	slices.SortFunc(input, func(a, b interface{}) bool { return checkOrder(a, b) > 0 })
	var k1, k2 int
	for i, v := range input {
		if string(util.Must(json.Marshal(v))) == "[[2]]" {
			k1 = i + 1
		}
		if string(util.Must(json.Marshal(v))) == "[[6]]" {
			k2 = i + 1
		}
	}
	return k1 * k2
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input))
	fmt.Println("Part2:", part2(input))
}
