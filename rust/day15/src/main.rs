#![allow(unused)]

use point::Point;

use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Vec<Reading>;
type Out = usize;
const PART1_RESULT: Out = 0;
const PART2_RESULT: Out = 0;

#[derive(Debug)]
struct Reading {
    sensor: Point,
    beacon: Point,
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    BufReader::new(input)
        .lines()
        .map(|l| {
            let l = l?;
            let s = l.split([' ', '=', ':', ',']).collect::<Vec<_>>();
            Ok(Reading {
                sensor: Point::new(
                    s[3].parse::<i32>().expect("i32"),
                    s[6].parse::<i32>().expect("i32"),
                ),
                beacon: Point::new(
                    s[13].parse::<i32>().expect("i32"),
                    s[16].parse::<i32>().expect("i32"),
                ),
            })
        })
        .collect::<std::io::Result<In>>()
}

fn project(r: &Reading, y_target: i32) -> Option<[i32; 2]> {
    let dx = r.sensor.manhattan(&r.beacon) - (r.sensor.y - y_target).abs();
    if dx > 0 {
        Some([r.sensor.x - dx, r.sensor.x + dx])
    } else {
        None
    }
}

fn part1(input: &In, y_target: i32) -> Out {
    for r in input {
        println!("{:?} -> {:?}", r, project(r, y_target));
    }
    PART1_RESULT
}

fn part2(input: &In) -> Out {
    PART2_RESULT
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&input, 2000000));
    println!("Part2: {:?}", part2(&input));
    Ok(())
}

#[test]
fn test_part1() {
    let input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&input, 10), PART1_RESULT);
}

#[test]
fn test_part2() {
    let input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&input), PART2_RESULT);
}

#[cfg(test)]
const TESTDATA: &str = "
Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3
";
