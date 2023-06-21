#![allow(unused)]

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

type In = Vec<Opcode>;
type Out = i32;
type Out2 = String;
const PART1_RESULT: Out = 13140;
const PART2_RESULT: &str = "\
##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....
";

#[derive(Debug)]
enum Opcode {
    NOOP,
    ADDX(i32),
}

impl TryFrom<&str> for Opcode {
    type Error = std::io::Error;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        match s {
            "noop" => Ok(Opcode::NOOP),
            _ => match s.split_once(" ") {
                Some(("addx", v)) => Ok(Opcode::ADDX(
                    v.parse::<i32>().map_err(|e| Error::new(InvalidData, s))?,
                )),
                _ => Err(Error::new(InvalidData, s)),
            },
        }
    }
}

#[derive(Debug)]
struct CRT([bool; 240]);

impl CRT {
    fn new() -> CRT {
        CRT([false; 240])
    }
}

impl std::fmt::Display for CRT {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for i in (0..240).step_by(40) {
            writeln!(
                f,
                "{}",
                &self.0[i..(i + 40)]
                    .iter()
                    .map(|&b| if b { "#" } else { "." })
                    .collect::<Vec<_>>()
                    .join("")
            )?;
        }
        Ok(())
    }
}

#[derive(Debug)]
struct CPU {
    x: i32,
    cycle: i32,
}

impl CPU {
    fn new() -> CPU {
        CPU { x: 1, cycle: 0 }
    }
    fn run<F>(&mut self, i: &Opcode, mut trace: F)
    where
        F: FnMut(i32, i32),
    {
        match i {
            Opcode::NOOP => {
                self.cycle += 1;
                trace(self.cycle, self.x);
            }
            Opcode::ADDX(n) => {
                for _ in 0..2 {
                    self.cycle += 1;
                    trace(self.cycle, self.x);
                }
                self.x += n;
            }
        }
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    BufReader::new(input)
        .lines()
        .map(|l| Opcode::try_from(l?.as_str()))
        .collect::<std::io::Result<In>>()
}

fn part1(input: &In) -> Out {
    let mut result = 0;
    let mut cpu = CPU::new();
    let mut trace_f = |cycle: i32, x: i32| {
        if (cycle - 20) % 40 == 0 {
            result += cycle * x;
        }
    };
    for i in input {
        cpu.run(i, &mut trace_f);
    }
    result
}

fn part2(input: &In) -> Out2 {
    let mut cpu = CPU::new();
    let mut crt = CRT::new();
    let mut trace_f = |cycle: i32, x: i32| {
        let pos = cycle - 1;
        if (pos % 40) >= x - 1 && (pos % 40) <= x + 1 {
            crt.0[pos as usize] = true;
        }
    };
    for i in input {
        cpu.run(i, &mut trace_f);
    }
    crt.to_string()
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {}", part1(&input));
    println!("Part2:\n{}", part2(&input));
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
addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop
";
