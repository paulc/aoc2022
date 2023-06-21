#![allow(unused)]

use std::collections::HashMap;
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
enum Move {
    Right(usize),
    Left(usize),
    Up(usize),
    Down(usize),
}

impl TryFrom<&str> for Move {
    type Error = std::io::Error;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        fn parse_usize(s: &str) -> Result<usize, std::io::Error> {
            s.parse::<usize>().map_err(|e| Error::new(InvalidData, s))
        }
        match s.split_once(" ") {
            Some(("R", s)) => Ok(Move::Right(parse_usize(s)?)),
            Some(("L", s)) => Ok(Move::Left(parse_usize(s)?)),
            Some(("U", s)) => Ok(Move::Up(parse_usize(s)?)),
            Some(("D", s)) => Ok(Move::Down(parse_usize(s)?)),
            _ => Err(Error::new(InvalidData, s)),
        }
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let reader = BufReader::new(input);
    reader
        .lines()
        .map(|l| Move::try_from(l?.as_str()))
        .collect::<std::io::Result<Vec<Move>>>()
}

fn part1(input: &In) -> Out {
    println!("{:?}", input);
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
