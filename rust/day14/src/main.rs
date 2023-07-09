#![allow(unused)]

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

mod grid;
mod point;

use crate::grid::Grid;
use crate::point::Offset;
use crate::point::Point;

type In = (Grid<Cave>, Grid<Cave>);
type Out = usize;
const PART1_RESULT: Out = 24;
const PART2_RESULT: Out = 93;

#[derive(Debug, Default, Clone)]
enum Cave {
    #[default]
    Air,
    Rock,
    Source,
    Sand,
}

impl Display for Cave {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                Cave::Air => ".",
                Cave::Rock => "#",
                Cave::Source => "+",
                Cave::Sand => "o",
            }
        )
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let paths = BufReader::new(input)
        .lines()
        .map(|l| {
            l.expect("Line")
                .split(" -> ")
                .map(|p| Point::try_from(p).expect("Invalid Point"))
                .collect::<Vec<_>>()
        })
        .collect::<Vec<_>>();
    let max_y = paths.iter().flatten().map(|p| p.y).max().unwrap();
    let min_x = paths.iter().flatten().map(|p| p.x).min().unwrap() - max_y;
    let max_x = paths.iter().flatten().map(|p| p.x).max().unwrap() + max_y;
    let mut c1: Grid<Cave> = Grid::new(Point::new(min_x, 0), Point::new(max_x, max_y));
    let mut c2: Grid<Cave> = Grid::new(Point::new(min_x, 0), Point::new(max_x, max_y + 2));
    c1.set(Point::new(500, 0), Cave::Source);
    c2.set(Point::new(500, 0), Cave::Source);
    paths.iter().for_each(|p| {
        p.windows(2).for_each(|p| {
            c1.draw_line(p[0], p[1], Cave::Rock);
            c2.draw_line(p[0], p[1], Cave::Rock);
        })
    });
    c2.draw_line(
        Point::new(min_x, max_y + 2),
        Point::new(max_x, max_y + 2),
        Cave::Rock,
    );
    Ok((c1, c2))
}

const DROP: [Offset; 3] = [
    Offset { dx: 0, dy: 1 },
    Offset { dx: -1, dy: 1 },
    Offset { dx: 1, dy: 1 },
];

fn drop_sand(cave: &mut Grid<Cave>) -> Option<()> {
    let mut p = Point::new(500, 0);
    let start = p.clone();
    let mut dropped = false;
    loop {
        dropped = false;
        for d in DROP {
            match cave.get(p + d) {
                Some(Cave::Air) => {
                    p = p + d;
                    dropped = true;
                    break;
                }
                None => return None,
                _ => {}
            }
        }
        if p == start {
            return None;
        }
        if !dropped {
            cave.set(p, Cave::Sand);
            return Some(());
        }
    }
}

fn part1(input: &mut In) -> Out {
    let cave = &mut input.0;
    let mut i: usize = 0;
    loop {
        match drop_sand(cave) {
            Some(()) => i += 1,
            None => break,
        }
    }
    i
}

fn part2(input: &mut In) -> Out {
    let cave = &mut input.1;
    let mut i: usize = 0;
    loop {
        match drop_sand(cave) {
            Some(()) => i += 1,
            None => break,
        }
    }
    i + 1
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let mut input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&mut input));
    println!("Part2: {:?}", part2(&mut input));
    Ok(())
}

#[test]
fn test_part1() {
    let mut input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&mut input), PART1_RESULT);
}

#[test]
fn test_part2() {
    let mut input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&mut input), PART2_RESULT);
}

#[cfg(test)]
const TESTDATA: &str = "
498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9
";
