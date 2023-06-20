#![allow(unused)]

use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Array2D<Tree>;
type Out = usize;
const PART1_RESULT: Out = 21;
const PART2_RESULT: Out = 8;

#[derive(Debug)]
struct ParseError(String);
impl std::error::Error for ParseError {}
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.0)
    }
}

#[derive(Debug, Clone)]
struct Tree(i8);

impl std::fmt::Display for Tree {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{} ", self.0)
    }
}

#[derive(Debug, Clone, Copy)]
enum Direction {
    North,
    South,
    East,
    West,
}

impl Direction {
    fn all() -> [Direction; 4] {
        [
            Direction::North,
            Direction::South,
            Direction::East,
            Direction::West,
        ]
    }
}

#[derive(Debug, Clone)]
struct Array2D<T> {
    data: Vec<Vec<T>>,
}

impl<T> Array2D<T> {
    fn new(data: Vec<Vec<T>>) -> Self {
        let (x_min, x_max, y_min, y_max) = (0, data[0].len(), 0, data.len());
        Self { data }
    }
    fn get(&self, (x, y): (usize, usize)) -> Option<&T> {
        self.data.get(y).and_then(|r| r.get(x))
    }
    fn get_mut(&mut self, (x, y): (usize, usize)) -> Option<&mut T> {
        self.data.get_mut(y).and_then(|r| r.get_mut(x))
    }
    fn iter(&self) -> impl Iterator<Item = (&T, usize, usize)> {
        self.data
            .iter()
            .enumerate()
            .flat_map(|(y, r)| r.iter().enumerate().map(move |(x, c)| (c, x, y)))
    }
    fn iter_mut(&mut self) -> impl Iterator<Item = (&mut T, usize, usize)> {
        self.data
            .iter_mut()
            .enumerate()
            .flat_map(|(y, r)| r.iter_mut().enumerate().map(move |(x, c)| (c, x, y)))
    }
    fn iter_direction(&self, (x, y): (usize, usize), direction: Direction) -> IterDirection<T> {
        IterDirection {
            array: self,
            direction,
            ix: x,
            iy: y,
        }
    }
}

impl<T: std::fmt::Display> std::fmt::Display for Array2D<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for y in 0..self.data.len() {
            writeln!(
                f,
                "{}",
                self.data[y]
                    .iter()
                    .map(|e| e.to_string())
                    .collect::<Vec<String>>()
                    .join("")
            )?
        }
        Ok(())
    }
}

struct IterDirection<'a, T> {
    array: &'a Array2D<T>,
    direction: Direction,
    ix: usize,
    iy: usize,
}

impl<'a, T> Iterator for IterDirection<'a, T> {
    type Item = (&'a T, usize, usize);
    fn next(&mut self) -> Option<Self::Item> {
        match self.direction {
            Direction::North => {
                if self.iy == 0 {
                    return None;
                } else {
                    self.iy -= 1
                }
            }
            Direction::South => self.iy += 1,
            Direction::East => self.ix += 1,
            Direction::West => {
                if self.ix == 0 {
                    return None;
                } else {
                    self.ix -= 1
                }
            }
        }
        self.array
            .data
            .get(self.iy)
            .and_then(|r| r.get(self.ix).and_then(|c| Some((c, self.ix, self.iy))))
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let reader = BufReader::new(input);
    Ok(Array2D::new(
        reader
            .lines()
            .map(|l| {
                l.and_then(|l| {
                    l.into_bytes()
                        .iter()
                        .map(|b| match b {
                            b'0'..=b'9' => Ok(Tree((b - b'0') as i8)),
                            _ => Err(Error::new(
                                InvalidData,
                                ParseError(format!("Invalid digit: {}", char::from(*b))),
                            )),
                        })
                        .collect::<Result<Vec<Tree>, Error>>()
                })
            })
            .collect::<Result<Vec<Vec<Tree>>, Error>>()?,
    ))
}

fn part1(input: &In) -> Out {
    input
        .iter()
        .map(|(t, x, y)| {
            // For each tree
            Direction::all()
                .iter() // For each direction
                .map(|&d| {
                    input
                        .iter_direction((x, y), d)
                        .map(|(t, _, _)| t.0)
                        .max() // Get the max tree height
                        .map(|max| t.0 > max) // Check if we can see over
                        .unwrap_or(true) // If None we are at edge do return true
                })
                .collect::<Vec<bool>>()
                .iter()
                .any(|&v| v) // Visible if any direction true
        })
        .filter(|&v| v) // Filter visible and count
        .count()
}

fn part2(input: &In) -> Out {
    let mut max: usize = 0;
    for (t, tx, ty) in input.iter() {
        // For each tree
        let mut view = [0, 0, 0, 0];
        for direction in Direction::all() {
            // For each direction
            for (v, vx, vy) in input.iter_direction((tx, ty), direction) {
                view[direction as usize] += 1; // Increment view length
                if t.0 <= v.0 {
                    // if tree height greater
                    break;
                }
            }
            // If any of the views is zero we can break
            if view[direction as usize] == 0 {
                break;
            }
        }
        let view = view.iter().fold(1_usize, |acc, &v| acc * v);
        if view > max {
            max = view;
        }
    }
    max
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
30373
25512
65332
33549
35390
";
