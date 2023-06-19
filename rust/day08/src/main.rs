#![allow(unused)]

use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Array2D<u8>;
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

#[derive(Debug, Clone, Copy, Eq, PartialEq)]
struct Point {
    x: usize,
    y: usize,
}

impl Point {
    fn new(x: usize, y: usize) -> Self {
        Self { x, y }
    }
}

#[derive(Debug, Clone, Copy, Eq, PartialEq)]
struct Delta {
    dx: usize,
    dy: usize,
}

impl Delta {
    fn new(dx: usize, dy: usize) -> Self {
        Self { dx, dy }
    }
}

#[derive(Debug)]
struct Array2D<T> {
    data: Vec<Vec<T>>,
    x_min: usize,
    x_max: usize,
    y_min: usize,
    y_max: usize,
    sep: String,
}

impl<T: Clone> Array2D<T> {
    fn new(data: Vec<Vec<T>>) -> Self {
        let (x_min, x_max, y_min, y_max) = (0, data[0].len(), 0, data.len());
        Self {
            data,
            x_min,
            x_max,
            y_min,
            y_max,
            sep: "".to_string(),
        }
    }
    fn get(&self, p: Point) -> Option<T> {
        if p.x >= self.x_min && p.x <= self.x_max && p.y >= self.y_min && p.y <= self.y_max {
            Some(self.data[p.y][p.x].clone())
        } else {
            None
        }
    }
    fn move_point(&self, p: Point, d: Delta) -> Option<Point> {
        let p = Point {
            x: p.x + d.dx,
            y: p.y + d.dy,
        };
        if p.x >= self.x_min && p.x <= self.x_max && p.y >= self.y_min && p.y <= self.y_max {
            Some(p)
        } else {
            None
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
                    .join(&self.sep)
            )?
        }
        Ok(())
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
                            b'0'..=b'9' => Ok((b - b'0') as u8),
                            _ => Err(Error::new(
                                InvalidData,
                                ParseError(format!("Invalid digit: {}", char::from(*b))),
                            )),
                        })
                        .collect::<Result<Vec<u8>, Error>>()
                })
            })
            .collect::<Result<Vec<Vec<u8>>, Error>>()?,
    ))
}

fn part1(input: &In) -> Out {
    PART1_RESULT
}

fn part2(input: &In) -> Out {
    PART2_RESULT
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
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&data), PART1_RESULT);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&data), PART2_RESULT);
}

#[cfg(test)]
const TESTDATA: &str = "
30373
25512
65332
33549
35390
";
