
use std::io::BufReader;
use std::io::prelude::*;
use std::fs::File;

fn main() -> std::io::Result<()> {
    let f = BufReader::new(File::open("input")?);
    let mut current: i32 = 0;
    let mut total: Vec<i32> = Vec::new();
    for l in f.lines() {
        let l = l?;
        if l.trim().is_empty() {
            total.push(current);
            current = 0;
        } else {
            current += l.trim().parse::<i32>().unwrap();
        }
    }
    if current > 0 {
        total.push(current);
    }
    total.sort();
    println!("Part 1: {}",total.iter().rev().take(1).sum::<i32>());
    println!("Part 2: {}",total.iter().rev().take(3).sum::<i32>());
    Ok(())
}
