#![allow(unused)]

mod packet;

use std::cmp::Ordering;
use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

use crate::packet::{Packet, PacketBuf, PacketRef};

type In = Vec<(PacketBuf<i32>, PacketBuf<i32>)>;
type Out = usize;
const PART1_RESULT: Out = 13;
const PART2_RESULT: Out = 140;

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut i = BufReader::new(input).lines();
    let mut out: Vec<(PacketBuf<i32>, PacketBuf<i32>)> = Vec::new();
    loop {
        let a = i.next().unwrap()?;
        let b = i.next().unwrap()?;
        let pa: PacketBuf<i32> = PacketBuf::try_from(a.as_str()).unwrap();
        let pb: PacketBuf<i32> = PacketBuf::try_from(b.as_str()).unwrap();
        out.push((pa, pb));
        match i.next() {
            None => break,
            _ => {}
        }
    }
    Ok(out)
}

fn value_to_list<T>(v: T) -> (PacketBuf<T>, PacketRef) {
    let mut buf: PacketBuf<T> = PacketBuf::new();
    let mut root = buf.get_ref(Packet::new_list(None));
    let p = buf.get_ref(Packet::Value(v));
    buf.push(root, p);
    (buf, root)
}

fn check_order(ba: &PacketBuf<i32>, bb: &PacketBuf<i32>, ia: PacketRef, ib: PacketRef) -> Ordering {
    // println!("Compare: {} - {}", ba.fmt_packet(ia), bb.fmt_packet(ib));
    let a = ba.0.get(ia.0).unwrap();
    let b = bb.0.get(ib.0).unwrap();
    match (a, b) {
        (Packet::Value(a), Packet::Value(b)) => a.cmp(b),
        (Packet::List(a, _), Packet::List(b, _)) => {
            let mut order = Ordering::Equal;
            let mut ia = a.iter();
            let mut ib = b.iter();
            loop {
                let a = ia.next();
                let b = ib.next();
                match (a, b) {
                    (Some(_), None) => order = Ordering::Greater,
                    (None, Some(_)) => order = Ordering::Less,
                    (Some(a), Some(b)) => order = check_order(ba, bb, *a, *b),
                    (None, None) => break,
                }
                if order != Ordering::Equal {
                    break;
                }
            }
            order
        }
        (Packet::Value(a), b) => {
            let (new_ba, new_a) = value_to_list(*a);
            check_order(&new_ba, bb, new_a, ib)
        }
        (a, Packet::Value(b)) => {
            let (new_bb, new_b) = value_to_list(*b);
            check_order(ba, &new_bb, ia, new_b)
        }
    }
}

fn part1(input: &In) -> Out {
    let mut result: usize = 0;
    for (i, (a, b)) in input.iter().enumerate() {
        if check_order(a, b, PacketRef(0), PacketRef(0)) == Ordering::Less {
            result += (i + 1);
        }
    }
    result
}

fn part2(input: &In) -> Out {
    PART2_RESULT
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&input));
    println!("Part2: {:?}", part2(&input));
    Ok(())
}

#[test]
fn test_part1() {
    let mut b = TESTDATA.trim_matches('\n').as_bytes();
    let input = parse_input(&mut b).unwrap();
    assert_eq!(part1(&input), PART1_RESULT);
}

#[test]
fn test_part2() {
    let mut b = TESTDATA.trim_matches('\n').as_bytes();
    let input = parse_input(&mut b).unwrap();
    assert_eq!(part2(&input), PART2_RESULT);
}

#[cfg(test)]
const TESTDATA: &str = "
[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]
";
