
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

impl Range {
    fn contains(&self,other: &Range) -> bool {
        (self.0 <= other.0) && (self.1 >= other.1)
    }
    fn overlaps(&self,other: &Range) -> bool {
        (self.1 >= other.0) && (self.0 <= other.1)
    }
}

#[derive(Debug)]
struct Pair(Range,Range);

fn parse_input(input: &mut impl Read) -> Vec<Pair> {
    let mut out: Vec<Pair> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            let (e1,e2 ) = l.split_once(',').unwrap();
            let r1 = Range::try_from(e1).unwrap();
            let r2 = Range::try_from(e2).unwrap();
            out.push(Pair(r1,r2));
        }
    }
    out
}

fn part1(input: &Vec<Pair>) -> Option<i32> {
    let mut score = 0;
    for p in input {
        if p.0.contains(&p.1) || p.1.contains(&p.0) {
            score += 1;
        }
    }
    Some(score)
}

fn part2(input: &Vec<Pair>) -> Option<i32> {
    let mut score = 0;
    for p in input {
        if p.0.overlaps(&p.1) || p.1.overlaps(&p.0) {
            score += 1;
        }
    }
    Some(score)
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
2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), 2);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), 4);
}

