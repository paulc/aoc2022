#![allow(unused)]

mod grid;
mod xy;

use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap, HashSet};
use std::fmt::Display;
use std::fs::File;
use std::hash::Hash;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

use crate::grid::Grid;
use crate::xy::XY;
use graph::dijkstra::shortest_path;
use graph::dijkstra::shortest_path_all;
use graph::Graph;

type In = Hill;
type Out = f64;
const PART1_RESULT: Out = 31.0;
const PART2_RESULT: Out = 29.0;

#[derive(Debug)]
struct Hill {
    hill: Grid<u8>,
    reachable: Graph<XY>,
    inverted: Graph<XY>,
    start: XY,
    end: XY,
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut start: XY = XY::new(0, 0);
    let mut end: XY = XY::new(0, 0);
    let mut hill = Grid::new(
        BufReader::new(input)
            .lines()
            .enumerate()
            .map(|(y, l)| {
                l?.chars()
                    .enumerate()
                    .map(|(x, c)| match c {
                        'a'..='z' => Ok((c as u8) - b'a'),
                        'S' => {
                            start = XY::new(x as i32, y as i32);
                            Ok(0)
                        }
                        'E' => {
                            end = XY::new(x as i32, y as i32);
                            Ok(25)
                        }
                        _ => Err(Error::new(InvalidData, format!("({},{}) {}", x, y, c))),
                    })
                    .collect::<std::io::Result<Vec<_>>>()
            })
            .collect::<std::io::Result<Vec<_>>>()?,
    );
    let mut reachable: Graph<XY> = Graph::new();
    for (h, p) in hill.iter() {
        hill.adjacent(p)
            .iter()
            .filter(|(&h1, _)| h1 <= h + 1)
            .for_each(|(_, p1)| {
                reachable.add_edge(p, *p1, 1.0);
            });
    }
    let mut inverted: Graph<XY> = Graph::new();
    for (h, p) in hill.iter() {
        hill.adjacent(p)
            .iter()
            .filter(|(&h1, _)| h1 >= if *h > 0 { *h - 1 } else { *h })
            .for_each(|(_, p1)| {
                inverted.add_edge(p, *p1, 1.0);
            });
    }

    Ok(Hill {
        hill,
        reachable,
        inverted,
        start,
        end,
    })
}

fn part1(input: &In) -> Out {
    let cost = shortest_path(&input.reachable, input.start, input.end)
        .unwrap()
        .0;
    cost
}

fn part2(input: &In) -> Out {
    let start = input
        .hill
        .iter()
        .filter(|(&h, _)| h == 0)
        .map(|(_, p)| p)
        .collect::<HashSet<_>>();
    let costs = shortest_path_all(&input.inverted, &input.end);
    let min = costs
        .iter()
        .filter(|(p, _)| start.contains(p))
        .map(|(_, h)| *h as i32)
        .min();
    min.unwrap() as f64
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&input));
    println!("Part2: {:?}", part2(&input));
    Ok(())
}

#[test]
fn test_part1() {
    let input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&input), PART1_RESULT);
}

#[test]
fn test_part2() {
    let input = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&input), PART2_RESULT);
}

#[cfg(test)]
const TESTDATA: &str = "
Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi
";
