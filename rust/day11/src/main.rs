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

type In = Vec<Monkey>;
type Out = usize;
const PART1_RESULT: Out = 10605;
const PART2_RESULT: Out = 2713310158;

#[derive(Debug, Clone)]
struct Monkey {
    inspected: usize,
    items: Vec<i64>,
    op: Option<Op>,
    test: Option<Test>,
    throw: (usize, usize),
}

#[derive(Debug, Copy, Clone)]
enum Op {
    Plus(i64),
    Mul(i64),
    Sqr,
}

#[derive(Debug, Copy, Clone)]
enum Test {
    Div(i64),
}

impl TryFrom<&str> for Monkey {
    type Error = std::io::Error;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        let mut m = Monkey {
            inspected: 0,
            items: Vec::new(),
            op: None,
            test: None,
            throw: (0, 0),
        };
        fn parse_i64(s: &str) -> std::io::Result<i64> {
            s.parse::<i64>().map_err(|e| Error::new(InvalidData, s))
        }
        fn parse_throw(s: &str) -> std::io::Result<usize> {
            match s.split_whitespace().last() {
                Some(s) => s.parse::<usize>().map_err(|e| Error::new(InvalidData, s)),
                None => Err(Error::new(InvalidData, s)),
            }
        }
        for l in s.split("\n") {
            if l.starts_with("Monkey") || l.is_empty() {
                continue;
            }
            match l.trim().split_once(": ") {
                Some(("Starting items", i)) => {
                    m.items = i
                        .split(", ")
                        .map(|s| parse_i64(s))
                        .collect::<Result<Vec<i64>, _>>()?;
                }
                Some(("Operation", op)) => {
                    match op.split_whitespace().collect::<Vec<_>>().as_slice() {
                        ["new", "=", "old", "*", "old"] => m.op = Some(Op::Sqr),
                        ["new", "=", "old", "*", n] => m.op = Some(Op::Mul(parse_i64(n)?)),
                        ["new", "=", "old", "+", n] => m.op = Some(Op::Plus(parse_i64(n)?)),
                        _ => return Err(Error::new(InvalidData, l)),
                    }
                }
                Some(("Test", t)) => match t.split_whitespace().collect::<Vec<_>>().as_slice() {
                    ["divisible", "by", n] => m.test = Some(Test::Div(parse_i64(n)?)),
                    _ => return Err(Error::new(InvalidData, l)),
                },
                Some(("If true", s)) => m.throw.0 = parse_throw(s)?,
                Some(("If false", s)) => m.throw.1 = parse_throw(s)?,
                _ => return Err(Error::new(InvalidData, l)),
            }
        }
        Ok(m)
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut reader = BufReader::new(input);
    let mut s = String::new();
    reader.read_to_string(&mut s)?;
    s.split("\n\n")
        .map(|s| Monkey::try_from(s))
        .collect::<std::io::Result<Vec<_>>>()
}

fn shuffle(input: &mut In, div: i64, lcm: i64) {
    for i in 0..input.len() {
        let mut throw: Vec<(i64, usize)> = Vec::new();
        let mut monkey = &mut input[i];
        while let Some(i) = monkey.items.pop() {
            let new = (match monkey.op.unwrap() {
                Op::Sqr => i * i,
                Op::Mul(n) => i * n,
                Op::Plus(n) => i + n,
            } / div)
                % lcm;
            let dest = match monkey.test.unwrap() {
                Test::Div(n) => {
                    if new % n == 0 {
                        monkey.throw.0
                    } else {
                        monkey.throw.1
                    }
                }
            };
            throw.push((new, dest));
            monkey.inspected += 1;
        }
        for (new, dest) in throw {
            input[dest].items.push(new);
        }
    }
}

fn part1(input: &mut In) -> Out {
    for round in 0..20 {
        shuffle(input, 3, i64::MAX);
    }
    input.sort_by_key(|m| m.inspected);
    input[input.len() - 2].inspected * input[input.len() - 1].inspected
}

fn part2(input: &mut In) -> Out {
    let mut lcm: i64 = 1;
    for m in input.iter() {
        lcm *= match m.test.unwrap() {
            Test::Div(n) => n,
        }
    }
    for round in 0..10000 {
        shuffle(input, 1, lcm);
    }
    input.sort_by_key(|m| m.inspected);
    input[input.len() - 2].inspected * input[input.len() - 1].inspected
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let mut input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&mut input.clone()));
    println!("Part2: {:?}", part2(&mut input));
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
Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1
";
