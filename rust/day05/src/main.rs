
#![allow(unused)]

use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

#[derive(Debug,Clone)]
struct Range(u32,u32);

impl TryFrom<&str> for Range {
    type Error = ();
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        match s.split_once('-') {
            Some((s1,s2)) => {
                match (s1.parse::<u32>(),s2.parse::<u32>()) {
                    (Ok(i1),Ok(i2)) => Ok(Self(i1,i2)),
                    _ => Err(())
                }
            }
            None => Err(()),
        }
    }
}

fn parse_input(input: &mut impl Read) -> Vec<Range> {
    let mut out: Vec<Range> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
        }
    }
    out
}

fn part1(input: &Vec<Range>) -> Option<&str> {
    Some("")
}

fn part2(input: &Vec<Range>) -> Option<&str> {
    Some("")
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let data = parse_input(&mut f);
    println!("Part1: {:?}",part1(&data).unwrap());
    println!("Part2: {:?}",part2(&data).unwrap());
    Ok(())
}

#[cfg(test)]
const TESTDATA : &str = "
    [D]
[N] [C]
[Z] [M] [P]
 1   2   3

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), "CMZ");
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), "MCD");
}

