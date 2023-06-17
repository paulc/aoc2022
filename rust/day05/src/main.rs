#![allow(unused)]
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

#[derive(Debug)]
struct Input {
    board: Vec<Vec<char>>,
    moves: Vec<Move>,
}

#[derive(Debug)]
struct Line(Vec<Option<char>>);

impl TryFrom<String> for Line {
    type Error = ();
    fn try_from(s: String) -> Result<Self, Self::Error> {
        Ok(Line(
            s.chars()
                .skip(1)
                .step_by(4)
                .map(|c| {
                    if c.is_ascii_uppercase() {
                        Some(c)
                    } else {
                        None
                    }
                })
                .collect::<Vec<Option<char>>>(),
        ))
    }
}

#[derive(Debug)]
struct Move {
    from: usize,
    to: usize,
    count: usize,
}

impl TryFrom<String> for Move {
    type Error = ();
    fn try_from(s: String) -> Result<Self, Self::Error> {
        let m: Vec<Result<usize, _>> = s
            .split_whitespace()
            .skip(1)
            .step_by(2)
            .map(|c| c.parse::<usize>())
            .collect();
        match *m.as_slice() {
            [Ok(count), Ok(from), Ok(to)] => Ok(Move { from, to, count }),
            _ => Err(()),
        }
    }
}

enum State {
    Stacks,
    Moves,
}

fn parse_input(input: &mut impl Read) -> Input {
    let mut out: Input = Input {
        board: Vec::new(),
        moves: Vec::new(),
    };
    let mut state = State::Stacks;
    let mut width = 0;
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            match state {
                State::Stacks => {
                    if l.len() == 0 {
                        state = State::Moves;
                    } else {
                        let line = Line::try_from(l).unwrap();
                        if width == 0 {
                            width = line.0.len();
                            for i in 0..width {
                                out.board.push(Vec::new());
                            }
                        }
                        for i in 0..width {
                            if let Some(c) = line.0[i] {
                                out.board[i].push(c);
                            }
                        }
                    }
                }
                State::Moves => {
                    out.moves.push(Move::try_from(l).unwrap());
                }
            }
        }
    }
    for s in out.board.iter_mut() {
        s.reverse()
    }
    out
}

fn part1(input: &Input) -> Option<String> {
    let mut board = input.board.clone();
    for m in &input.moves {
        for i in 0..m.count {
            let c = board[m.from - 1].pop().unwrap();
            board[m.to - 1].push(c);
        }
    }
    let out: String = board.iter_mut().map(|r| r.pop().unwrap()).collect();
    Some(out)
}

fn part2(input: &Input) -> Option<String> {
    let mut board = input.board.clone();
    for m in &input.moves {
        let at = board[m.from - 1].len() - m.count;
        for c in board[m.from - 1].split_off(at) {
            board[m.to - 1].push(c);
        }
    }
    let out: String = board.iter_mut().map(|r| r.pop().unwrap()).collect();
    Some(out)
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let data = parse_input(&mut f);
    println!("Part1: {:?}", part1(&data).unwrap());
    println!("Part2: {:?}", part2(&data).unwrap());
    Ok(())
}

#[cfg(test)]
const TESTDATA: &str = "
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
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes());
    assert_eq!(part1(&data).unwrap(), "CMZ");
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes());
    assert_eq!(part2(&data).unwrap(), "MCD");
}
