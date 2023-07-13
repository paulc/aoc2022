#![allow(unused)]

mod rangelist;

use point::Point;
use rangelist::Rangelist;

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::env;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;
use std::sync::mpsc;
use std::thread;

type In = Vec<Reading>;
type Out = i64;
const PART1_RESULT: Out = 26;
const PART2_RESULT: Out = 56000011;

#[derive(Debug, Clone)]
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
    let mut beacons: HashSet<i32> = HashSet::new();
    let mut rl: Rangelist<i32> = Rangelist::new();
    input
        .iter()
        .filter_map(|r| {
            if r.beacon.y == y_target {
                beacons.insert(r.beacon.x);
            }
            project(r, y_target)
        })
        .for_each(|r| rl.add(r));
    rl.coalesce();
    let n_beacons = beacons.into_iter().filter(|&x| rl.contains(x)).count();
    ((rl.iter().map(|&[a, b]| b - a + 1).sum::<i32>() as usize) - n_beacons) as i64
}

fn part2(input: &In, max_xy: i32) -> Out {
    let nthread = match env::var_os("THREAD") {
        Some(_) => match thread::available_parallelism() {
            Ok(n) => n.get(),
            _ => 1,
        },
        None => 1,
    };
    let step = max_xy as usize / nthread;
    let mut result: i64 = 0;
    let (tx, rx) = mpsc::channel();
    for i in 0..nthread {
        let tx = tx.clone();
        let input = input.clone();
        thread::spawn(move || {
            for y in (step * i)..(step * i + step) {
                let mut rl: Rangelist<i32> = Rangelist::new();
                input
                    .iter()
                    .filter_map(|r| project(r, y.try_into().unwrap()))
                    .for_each(|r| rl.add(r));
                rl.coalesce();
                if rl.len() > 1 {
                    let l = rl.pop_first().unwrap();
                    let r = rl.pop_first().unwrap();
                    if r[0] > l[1] + 1 {
                        result = (l[1] as i64 + 1) * 4000000 + y as i64;
                        tx.send(result);
                    }
                }
            }
        });
    }
    rx.recv().unwrap()
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
