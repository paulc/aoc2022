
use std::fs::File;
use std::io::prelude::*;

fn parse_input(input: &mut impl Read) -> Result<Vec<Vec<i32>>, std::io::Error> {
    let mut s = String::new();
    input.read_to_string(&mut s)?;
    Ok(s.lines().collect::<Vec<_>>()
                    .split(|line| line.is_empty())
                    .map(|group| { group
                                    .iter()
                                    .map(|v| v.parse::<i32>().unwrap())
                                    .collect::<Vec<_>>()
                    })
                    .collect::<Vec<_>>()
    )
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
    let data = parse_input(&mut f)?;
    println!("Part1: {:?}",part1(&data).unwrap());
    println!("Part2: {:?}",part2(&data).unwrap());

    Ok(())
}

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
    let data = parse_input(&mut TESTDATA.trim().as_bytes()).unwrap();
    assert_eq!(part1(&data).unwrap(), 24000);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim().as_bytes()).unwrap();
    assert_eq!(part2(&data).unwrap(), 45000);
}

