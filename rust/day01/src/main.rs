
use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

fn parse_input(input: &mut impl Read) -> Vec<Vec<i32>> {
    let mut out: Vec<Vec<i32>> = Vec::new();
    let mut current: Vec<i32> = Vec::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            if l.is_empty() && !current.is_empty() {
                    out.push(current.clone());
                    current.clear();
            } else {
                if let Ok(i) = l.parse::<i32>() {
                    current.push(i);
                }
            }
        }
    }
    if !current.is_empty() {
        out.push(current);
    }
    out
}

fn part1(input: &Vec<Vec<i32>>) -> Option<i32> {
    input.iter().map(|group| group.iter().sum()).max()
}

fn part2(input: &Vec<Vec<i32>>) -> Option<i32> {
    let mut groups = input.iter().map(|group| group.iter().sum::<i32>()).collect::<Vec<_>>();
    groups.sort();
    Some(groups.iter().rev().take(3).sum::<i32>())
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
1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
";

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part1(&data).unwrap(), 24000);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes());
    assert_eq!(part2(&data).unwrap(), 45000);
}

