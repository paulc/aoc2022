#![allow(unused)]

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Vec<Move>;
type Out = usize;
const PART1_RESULT: Out = 13;
const PART2_RESULT: Out = 1;
const PART2_RESULT2: Out = 36;

#[derive(Debug)]
enum Direction {
    Right,
    Left,
    Up,
    Down,
}

#[derive(Debug)]
struct Move {
    d: Direction,
    n: usize,
}

impl TryFrom<&str> for Move {
    type Error = std::io::Error;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        let parse_usize = |s: &str| s.parse::<usize>().map_err(|e| Error::new(InvalidData, s));
        match s.split_once(" ") {
            Some(("R", s)) => Ok(Move {
                d: Direction::Right,
                n: parse_usize(s)?,
            }),
            Some(("L", s)) => Ok(Move {
                d: Direction::Left,
                n: parse_usize(s)?,
            }),
            Some(("U", s)) => Ok(Move {
                d: Direction::Up,
                n: parse_usize(s)?,
            }),
            Some(("D", s)) => Ok(Move {
                d: Direction::Down,
                n: parse_usize(s)?,
            }),
            _ => Err(Error::new(InvalidData, s)),
        }
    }
}

#[derive(Debug, Clone, Copy, Eq, PartialEq, Default, Hash)]
struct XY(i32, i32);

impl XY {
    fn is_touching(&self, other: &XY) -> bool {
        ((self.0 - other.0).abs() <= 1) && ((self.1 - other.1).abs() <= 1)
    }
    fn add(&self, other: &XY) -> XY {
        XY(self.0 + other.0, self.1 + other.1)
    }
    fn cmp(&self, other: &XY) -> (Ordering, Ordering) {
        (self.0.cmp(&(other.0)), self.1.cmp(&(other.1)))
    }
}

#[derive(Debug)]
struct Rope {
    knots: Vec<XY>,
    visited: HashSet<XY>,
}

impl Rope {
    fn new(n: usize) -> Self {
        Self {
            knots: vec![XY::default(); n],
            visited: HashSet::new(),
        }
    }
    fn move_rope(&mut self, m: &Move) {
        for _ in 0..m.n {
            // Move head
            match m.d {
                Direction::Right => self.knots[0] = self.knots[0].add(&XY(1, 0)),
                Direction::Left => self.knots[0] = self.knots[0].add(&XY(-1, 0)),
                Direction::Up => self.knots[0] = self.knots[0].add(&XY(0, 1)),
                Direction::Down => self.knots[0] = self.knots[0].add(&XY(0, -1)),
            }
            // Move tail
            for i in 1..self.knots.len() {
                if !self.knots[i].is_touching(&self.knots[i - 1]) {
                    let (cmp_x, cmp_y) = self.knots[i - 1].cmp(&self.knots[i]);
                    let offset = XY(
                        match cmp_x {
                            Ordering::Less => -1,
                            Ordering::Equal => 0,
                            Ordering::Greater => 1,
                        },
                        match cmp_y {
                            Ordering::Less => -1,
                            Ordering::Equal => 0,
                            Ordering::Greater => 1,
                        },
                    );
                    self.knots[i] = self.knots[i].add(&offset);
                }
            }
            // Mark where tail visits
            self.visited.insert(self.knots[self.knots.len() - 1]);
        }
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    BufReader::new(input)
        .lines()
        .map(|l| Move::try_from(l?.as_str()))
        .collect::<std::io::Result<Vec<Move>>>()
}

fn part1(input: &In) -> Out {
    let mut rope = Rope::new(2);
    input.iter().for_each(|m| rope.move_rope(m));
    rope.visited.len()
}

fn part2(input: &In) -> Out {
    let mut rope = Rope::new(10);
    input.iter().for_each(|m| rope.move_rope(m));
    rope.visited.len()
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
    let input = parse_input(&mut TESTDATA2.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&input), PART2_RESULT2);
}

#[cfg(test)]
const TESTDATA: &str = "
R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2
";

#[cfg(test)]
const TESTDATA2: &str = "
R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20
";
