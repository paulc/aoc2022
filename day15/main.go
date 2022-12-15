package main

import (
	"fmt"
	"io"
	"os"

	"github.com/paulc/aoc2022/util"
	"github.com/paulc/aoc2022/util/point"
	"github.com/paulc/aoc2022/util/reader"
	"github.com/paulc/aoc2022/util/set"
	"golang.org/x/exp/slices"
)

func parseInput(r io.Reader) [][2]point.Point {
	return util.Map(util.Must(reader.Lines(r)), func(s string) (out [2]point.Point) {
		//var i [4]int
		//util.Must(fmt.Sscanf(s, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &i[0], &i[1], &i[2], &i[3]))
		util.Must(fmt.Sscanf(s, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &out[0].X, &out[0].Y, &out[1].X, &out[1].Y))
		return
		// return [2]point.Point{point.Point{i[0], i[1]}, point.Point{i[2], i[3]}}
	})
}

func calculateExcluded(input [][2]point.Point, targetY int) (excluded [][2]int, beacons set.Set[int]) {
	temp := [][2]int{}
	beacons = set.NewSet[int]()
	for _, v := range input {
		if v[1].Y == targetY {
			beacons.Add(v[1].X)
		}
		target := point.Point{v[0].X, targetY}
		d := v[0].Distance(v[1])
		dx := d - v[0].Ydistance(target)
		if dx >= 0 {
			temp = append(temp, [2]int{v[0].X - dx, v[0].X + dx})
		}
	}
	slices.SortFunc(temp, func(a, b [2]int) bool {
		if a[0] == b[0] {
			return a[1] < b[1]
		}
		return a[0] < b[0]
	})
	for i, _ := range temp {
		if i < len(temp)-1 {
			if temp[i][1]+1 >= temp[i+1][0] {
				temp[i+1][0] = temp[i][0]
				temp[i+1][1] = util.Max(temp[i][1], temp[i+1][1])
			} else {
				excluded = append(excluded, temp[i])
			}
		} else {
			excluded = append(excluded, temp[i])
		}
	}
	return
}

func part1(input [][2]point.Point, targetY int) (result int) {
	excluded, beacons := calculateExcluded(input, targetY)
	for _, v := range excluded {
		result += v[1] - v[0] + 1
		for b := range beacons {
			if b >= v[0] && b <= v[1] {
				result--
			}
		}
	}
	return
}

func part2(input [][2]point.Point, maxXY int) (result int) {
	for i := 0; i < maxXY; i++ {
		excluded, _ := calculateExcluded(input, i)
		if len(excluded) > 1 {
			result = i + (excluded[0][1]+1)*4000000
		}
	}
	return
}

func main() {
	input := parseInput(util.Must(os.Open("input")))
	fmt.Println("Part1:", part1(input, 2000000))
	fmt.Println("Part2:", part2(input, 4000000))
}
