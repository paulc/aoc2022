#![allow(unused)]

use point::Point;

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Vec<Reading>;
type Out = i64;
const PART1_RESULT: Out = 26;
const PART2_RESULT: Out = 56000011;

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

fn coalesce(r: &Vec<[i32; 2]>) -> Vec<[i32; 2]> {
    let mut i = r.into_iter();
    let mut out = Vec::new();
    match i.next() {
        None => out,
        Some(start) => {
            let mut current = start.clone();
            while let Some(next) = i.next() {
                match (
                    current[0].cmp(&next[0]),
                    current[1].cmp(&next[0]),
                    current[1].cmp(&next[1]),
                ) {
                    // c0---c1 n0---n1
                    (Ordering::Less, Ordering::Less, _) => {
                        out.push(current);
                        current = next.clone();
                    }
                    // c0---c1
                    //    n0---n1
                    (
                        Ordering::Less | Ordering::Equal,
                        Ordering::Greater | Ordering::Equal,
                        Ordering::Less | Ordering::Equal,
                    ) => {
                        current = [current[0], next[1]];
                    }
                    // c0--------c1
                    //    n0---n1
                    (Ordering::Less | Ordering::Equal, Ordering::Greater, Ordering::Greater) => {
                        current = [current[0], current[1]];
                    }
                    _ => panic!("Shouldnt get here"),
                }
            }
            out.push(current);
            out
        }
    }
}

fn inside(cover: &Vec<[i32; 2]>, x: &i32) -> bool {
    for [a, b] in cover {
        if x >= a && x <= b {
            return true;
        }
    }
    false
}

fn part1(input: &In, y_target: i32) -> Out {
    let mut beacons: HashSet<i32> = HashSet::new();
    let mut cover = input
        .iter()
        .filter_map(|r| {
            if r.beacon.y == y_target {
                beacons.insert(r.beacon.x);
            }
            project(r, y_target)
        })
        .collect::<Vec<_>>();
    cover.sort();
    let cover = coalesce(&cover);
    let n_beacons = beacons.iter().filter(|&x| inside(&cover, x)).count();
    ((cover.iter().map(|&[a, b]| b - a + 1).sum::<i32>() as usize) - n_beacons) as i64
}

fn part2(input: &In, max_xy: i32) -> Out {
    let mut result: i64 = 0;
    for y in 0..=max_xy {
        let mut cover = input
            .iter()
            .filter_map(|r| project(r, y))
            .collect::<Vec<_>>();
        cover.sort();
        let cover = coalesce(&cover);
        if cover.len() > 1 {
            result = (cover[0][1] as i64 + 1) * 4000000 + y as i64;
            break;
        }
    }
    result
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&input, 2000000));
    println!("Part2: {:?}", part2(&input, 4000000));
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
    assert_eq!(part2(&input, 20), PART2_RESULT);
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
