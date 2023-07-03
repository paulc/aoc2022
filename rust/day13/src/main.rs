#![allow(unused)]

mod packet;

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Vec<()>;
type Out = usize;
const PART1_RESULT: Out = 0;
const PART2_RESULT: Out = 0;

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
    BufReader::new(input)
        .lines()
        .map(|l| Ok(()))
        .collect::<std::io::Result<In>>()
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
";
