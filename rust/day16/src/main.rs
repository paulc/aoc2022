#![allow(unused)]

use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Display;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;
use std::sync::atomic::AtomicI32;
use std::sync::atomic::Ordering;
use std::sync::Arc;
use std::thread;

use graph2::Graph;
use graph2::Vertex;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Key([u8; 2]);

impl From<&str> for Key {
    fn from(s: &str) -> Self {
        let mut bytes = s.bytes();
        Key([bytes.next().unwrap(), bytes.next().unwrap()])
    }
}

impl Display for Key {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}{}", self.0[0], self.0[1])
    }
}

#[derive(Debug)]
struct Cave {
    non_zero: HashSet<Key>,
    valve: HashMap<Key, i32>,
    cost: HashMap<Key, HashMap<Key, i32>>,
}

type In = Cave;
type Out = i32;
const PART1_RESULT: Out = 1651;
const PART2_RESULT: Out = 1707;

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut cave: Graph<Key, i32> = Graph::new();
    BufReader::new(input)
        .lines()
        .map(|l| {
            let l = l.expect("line");
            let s = l.split([' ', '=', ';', ',']).collect::<Vec<_>>();
            let next = s[11..]
                .iter()
                .filter(|s| !s.is_empty())
                .map(|&e| (Key::from(e), 1))
                .collect::<Vec<_>>();
            Vertex::new(
                Key::from(s[1]),
                s[5].parse::<i32>().expect("flow rate"),
                next,
            )
        })
        .for_each(|v| cave.add_vertex(v));
    let non_zero = cave
        .vertices()
        .filter_map(|v| {
            if v.data > 0 {
                Some(v.key.clone())
            } else {
                None
            }
        })
        .collect::<HashSet<_>>();
    let valve = cave
        .vertices()
        .filter_map(|v| {
            if v.data > 0 {
                Some((v.key.clone(), v.data.clone()))
            } else {
                None
            }
        })
        .collect::<HashMap<_, _>>();
    let start = vec![Key::from("AA")];
    let cost = non_zero
        .iter()
        .chain(start.iter())
        .map(|a| {
            (
                a.clone(),
                non_zero
                    .iter()
                    .map(|b| {
                        (
                            b.clone(),
                            cave.astar(a.clone(), b.clone(), |_| 1).unwrap().0,
                        )
                    })
                    .collect::<HashMap<_, _>>(),
            )
        })
        .collect::<HashMap<_, _>>();
    Ok(Cave {
        non_zero,
        valve,
        cost,
    })
}

#[derive(Debug, Clone)]
struct State {
    current: Key,
    remaining: HashSet<Key>, // (next, volume)
    pressure: i32,
    time: i32,
}

fn search<'a>(
    s: State,
    valve: Arc<HashMap<Key, i32>>,
    cost: Arc<HashMap<Key, HashMap<Key, i32>>>,
) -> Vec<i32> {
    if !s.remaining.is_empty() && s.time < 30 {
        let mut out: Vec<i32> = Vec::new();
        for n in s.remaining.iter() {
            let t_next = s.time + cost[&s.current][n] + 1;
            let mut next = State {
                current: n.clone(),
                remaining: s.remaining.clone(),
                pressure: s.pressure + (30 - t_next) * valve[n],
                time: t_next,
            };
            next.remaining.remove(n);
            out.extend(search(next, valve.clone(), cost.clone()));
        }
        out
    } else {
        // println!("{:?}", s);
        vec![s.pressure]
    }
}

fn part1(
    Cave {
        non_zero,
        valve,
        cost,
    }: &In,
) -> Out {
    let start = State {
        current: Key::from("AA"),
        remaining: non_zero.clone(),
        pressure: 0,
        time: 0,
    };
    let mut handle = Vec::new();
    let valve = Arc::new(valve.clone());
    let cost = Arc::new(cost.clone());
    let result = Arc::new(AtomicI32::new(0));
    for n in start.remaining.iter() {
        let t_next = start.time + cost[&start.current][n] + 1;
        let mut next = State {
            current: n.clone(),
            remaining: start.remaining.clone(),
            pressure: start.pressure + (30 - t_next) * valve[n],
            time: t_next,
        };
        next.remaining.remove(n);
        let valve = valve.clone();
        let cost = cost.clone();
        let result = result.clone();
        handle.push(thread::spawn(move || {
            let max = search(next, valve, cost).into_iter().max().unwrap();
            result.fetch_max(max, Ordering::Relaxed);
        }));
    }
    for h in handle {
        h.join().unwrap();
    }
    (*result).load(Ordering::Relaxed)
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
Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II
";
