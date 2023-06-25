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

type In = Hill;
type Out = usize;
const PART1_RESULT: Out = 0;
const PART2_RESULT: Out = 0;

#[derive(Debug)]
struct Hill {
    hill: Vec<Vec<u8>>,
    start: (usize, usize),
    end: (usize, usize),
}

/*
impl TryFrom<&str> for ____ {
    type Error = std::io::Error;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        Err(Error::new(InvalidData, "Error")),
    }
}

impl Display for ____ {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        writeln!(f,"{}",____)
    }
}
*/

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut start = (0, 0);
    let mut end = (0, 0);
    let mut hill = BufReader::new(input)
        .lines()
        .enumerate()
        .map(|(y, l)| {
            l?.chars()
                .enumerate()
                .map(|(x, c)| match c {
                    'a'..='z' => Ok((c as u8) - b'a'),
                    'S' => {
                        start = (x, y);
                        Ok(0)
                    }
                    'E' => {
                        end = (x, y);
                        Ok(25)
                    }
                    _ => Err(Error::new(InvalidData, format!("({},{}) {}", x, y, c))),
                })
                .collect::<std::io::Result<Vec<_>>>()
        })
        .collect::<std::io::Result<Vec<_>>>()?;
    dbg!(Ok(Hill { hill, start, end }))
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
