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

fn packet_cmp(
    buf_a: &PacketBuf<i32>,
    buf_b: &PacketBuf<i32>,
    index_a: PacketRef,
    index_b: PacketRef,
) -> Ordering {
    // println!("Compare: {} - {}", buf_a.fmt_packet(ia), buf_b.fmt_packet(ib));
    let a = buf_a.0.get(index_a.0).unwrap();
    let b = buf_b.0.get(index_b.0).unwrap();
    match (a, b) {
        (Packet::Value(a), Packet::Value(b)) => a.cmp(b),
        (Packet::List(a, _), Packet::List(b, _)) => {
            let mut order = Ordering::Equal;
            let mut iter_a = a.iter();
            let mut iter_b = b.iter();
            loop {
                order = match (iter_a.next(), iter_b.next()) {
                    (Some(_), None) => Ordering::Greater,
                    (None, Some(_)) => Ordering::Less,
                    (Some(a), Some(b)) => packet_cmp(buf_a, buf_b, *a, *b),
                    (None, None) => break,
                };
                if order != Ordering::Equal {
                    break;
                }
            }
            order
        }
        (Packet::Value(a), b) => {
            let (buf, root) = value_to_list(*a);
            packet_cmp(&buf, buf_b, root, index_b)
        }
        (a, Packet::Value(b)) => {
            let (buf, root) = value_to_list(*b);
            packet_cmp(buf_a, &buf, index_a, root)
        }
    }
}

fn part1(input: &In) -> Out {
    let mut result: usize = 0;
    for (i, (a, b)) in input.iter().enumerate() {
        if packet_cmp(a, b, PacketRef(0), PacketRef(0)) == Ordering::Less {
            result += (i + 1);
        }
    }
    result
}

fn part2(input: &In) -> Out {
    let mut flat = input
        .iter()
        .flat_map(|(a, b)| vec![a, b])
        .collect::<Vec<_>>();
    let f1 = &PacketBuf::try_from("[[2]]").unwrap();
    let f2 = &PacketBuf::try_from("[[6]]").unwrap();
    flat.push(f1);
    flat.push(f2);
    flat.sort_by(|a, b| packet_cmp(a, b, PacketRef(0), PacketRef(0)));
    let f1 = flat.iter().position(|&p| p.to_string() == "[[2]]").unwrap() + 1;
    let f2 = flat.iter().position(|&p| p.to_string() == "[[6]]").unwrap() + 1;
    f1 * f2
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
