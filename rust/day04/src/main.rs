
#![allow(unused)]

use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

#[derive(Debug,Clone,Copy,Eq,Hash,PartialEq)]
struct Item {
    item: u8,
    priority: i32,
}

impl TryFrom<u8> for Item {
    type Error = &'static str;
    fn try_from(b: u8) -> Result<Self, Self::Error> {
        match b {
            b'a'..=b'z' => Ok(Self{ item: b, priority: (b-b'a'+1) as i32 }),
            _   => Err("Invalid Item")
        }
    }
}

impl std::fmt::Display for Item {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", std::str::from_utf8(&[self.item]).unwrap())
    }
}

fn parse_input(input: &mut impl Read) -> Vec<Item> {
    let mut out: Vec<Item> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
        }
    }
    out
}

fn part1(input: &Vec<Item>) -> Option<i32> {
    let mut score = 0;
    Some(score)
}

fn part2(input: &Vec<Item>) -> Option<i32> {
    let mut score = 0;
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
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), 0);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), 0);
}

