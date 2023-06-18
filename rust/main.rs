#![allow(unused)]

use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = ();
type Out = ();
const PART1_RESULT: Out = ();
const PART2_RESULT: Out = ();

#[derive(Debug)]
struct ParseError;
impl std::error::Error for ParseError {}
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Parse error")
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    Ok(())
}

fn part1(input: &In) -> Out {
    ()
}

fn part2(input: &In) -> Out {
    ()
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
";
