#![allow(unused)]
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

fn parse_input(input: &mut impl Read) -> std::io::Result<Vec<u8>> {
    let mut out: Vec<u8> = Vec::new();
    input.read_to_end(&mut out)?;
    Ok(out)
}

fn find_start(input: &Vec<u8>, start_len: usize) -> Option<usize> {
    for i in 0..(input.len() - start_len) {
        let s: HashSet<u8> = input.iter().skip(i).take(start_len).cloned().collect();
        if s.len() == start_len {
            return Some(i + start_len);
        }
    }
    None
}

fn part1(input: &Vec<u8>) -> Option<usize> {
    find_start(input, 4)
}

fn part2(input: &Vec<u8>) -> Option<usize> {
    find_start(input, 14)
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let data = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&data).unwrap());
    println!("Part2: {:?}", part2(&data).unwrap());
    Ok(())
}

#[cfg(test)]
const TESTDATA: &str = "
mjqjpqmgbljsphdztnvjfqwrcgsmlb
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&data).unwrap(), 7);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&data).unwrap(), 19);
}
