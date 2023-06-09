
#![allow(unused)]

use std::fs::File;
use std::io::BufReader;
use std::collections::HashSet;
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
            b'A'..=b'Z' => Ok(Self{ item: b, priority: (b-b'A'+27) as i32 }),
            _   => Err("Invalid Item")
        }
    }
}

impl std::fmt::Display for Item {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", std::str::from_utf8(&[self.item]).unwrap())
    }
}

#[derive(Debug)]
struct Rucsac {
    all: HashSet<Item>,
    c1: HashSet<Item>,
    c2: HashSet<Item>,
}

impl Rucsac {
    fn new() -> Self {
        Self{ all: HashSet::new(), c1: HashSet::new(), c2: HashSet::new() }
    }
}

impl TryFrom<&str> for Rucsac {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        let mut out = Rucsac::new();
        for i in s.as_bytes() {
            out.all.insert(Item::try_from(*i)?);
        }
        let (s1,s2) = s.split_at(s.len()/2);
        for i in s1.as_bytes() {
            out.c1.insert(Item::try_from(*i)?);
        }
        for i in s2.as_bytes() {
            out.c2.insert(Item::try_from(*i)?);
        }
        Ok(out)
    }
}

impl std::fmt::Display for Rucsac {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let v1: Vec<_> = self.c1.iter().map(|item| item.to_string()).collect();
        let v2: Vec<_> = self.c2.iter().map(|item| item.to_string()).collect();
        write!(f, ">> {} :: {}", v1.join(""), v2.join(""))
    }
}

fn parse_input(input: &mut impl Read) -> Vec<Rucsac> {
    let mut out: Vec<Rucsac> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            out.push(Rucsac::try_from(l.as_str()).unwrap());
        }
    }
    out
}

fn part1(input: &Vec<Rucsac>) -> Option<i32> {
    let mut score = 0;
    for r in input {
        score += r.c1.intersection(&r.c2).fold(0, |acc,i| acc + i.priority);
    }
    Some(score)
}

fn part2(input: &Vec<Rucsac>) -> Option<i32> {
    let mut score = 0;
    for g in input.chunks(3) {
        let badge: HashSet<_> = g[0].all.intersection(&(g[1].all)).cloned().collect();
        score += badge.intersection(&(g[2].all)).fold(0, |acc,i| acc + i.priority);
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
vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), 157);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), 70);
}

