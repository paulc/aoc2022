
use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

#[derive(Debug)]
enum Move {
    Rock,
    Paper,
    Scissors,
}

impl TryFrom<&str> for Move {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        match s {
            "A"|"X" => Ok(Self::Rock),
            "B"|"Y" => Ok(Self::Paper),
            "C"|"Z" => Ok(Self::Scissors),
            _   => Err("Invalid Move")
        }
    }
}

#[derive(Debug)]
enum Outcome {
    Lose,
    Draw,
    Win,
}

impl TryFrom<&str> for Outcome {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        match s {
            "X" => Ok(Self::Lose),
            "Y" => Ok(Self::Draw),
            "Z" => Ok(Self::Win),
            _   => Err("Invalid Result")
        }
    }
}

#[derive(Debug)]
struct Turn {
    elf_move: Move,
    you_move: Move,
}

#[derive(Debug)]
struct Turn2 {
    elf_move: Move,
    result: Outcome,
}


fn parse_input(input: &mut impl Read) -> Vec<Turn> {
    let mut out: Vec<Turn> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            let split = l.splitn(2," ").collect::<Vec<_>>();
            out.push(Turn{ elf_move: Move::try_from(split[0]).unwrap(),
                           you_move: Move::try_from(split[1]).unwrap() });
        }
    }
    out
}

fn parse_input2(input: &mut impl Read) -> Vec<Turn2> {
    let mut out: Vec<Turn2> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            let split = l.splitn(2," ").collect::<Vec<_>>();
            out.push(Turn2{ elf_move: Move::try_from(split[0]).unwrap(),
                            result: Outcome::try_from(split[1]).unwrap() });
        }
    }
    out
}


fn part1(input: &Vec<Turn>) -> Option<i32> {
    let mut score = 0;
    for turn in input.iter() {
        score += match turn {
            Turn{elf_move: elf, you_move: Move::Rock} => 
                        match elf { Move::Rock => 4, Move::Paper => 1, Move::Scissors => 7 },
            Turn{elf_move: elf, you_move: Move::Paper} => 
                        match elf { Move::Rock => 8, Move::Paper => 5, Move::Scissors => 2 },
            Turn{elf_move: elf, you_move: Move::Scissors} => 
                        match elf { Move::Rock => 3, Move::Paper => 9, Move::Scissors => 6 },
        }
    }
    Some(score)
}

fn part2(input: &Vec<Turn2>) -> Option<i32> {
    let mut score = 0;
    for turn in input.iter() {
        score += match turn {
            Turn2{elf_move: elf, result: Outcome::Lose} => 
                        match elf { Move::Rock => 3, Move::Paper => 1, Move::Scissors => 2 },
            Turn2{elf_move: elf, result: Outcome::Draw} => 
                        match elf { Move::Rock => 4, Move::Paper => 5, Move::Scissors => 6 },
            Turn2{elf_move: elf, result: Outcome::Win} => 
                        match elf { Move::Rock => 8, Move::Paper => 9, Move::Scissors => 7 },
        }
    }
    Some(score)
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let data = parse_input(&mut f);
    println!("Part1: {:?}",part1(&data).unwrap());
    let mut f2 = File::open("input")?;
    let data2 = parse_input2(&mut f2);
    println!("Part2: {:?}",part2(&data2).unwrap());
    Ok(())
}

#[cfg(test)]
const TESTDATA : &str = "
A Y
B X
C Z
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), 15);
}

#[test]
fn test_part2() {
    let data = parse_input2(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), 12);
}

