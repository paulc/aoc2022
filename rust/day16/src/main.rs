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

#[derive(Clone, Copy, PartialEq, Eq, Hash, PartialOrd, Ord)]
struct Key([u8; 2]);

impl From<&str> for Key {
    fn from(s: &str) -> Self {
        let mut bytes = s.bytes();
        Key([bytes.next().unwrap(), bytes.next().unwrap()])
    }
}

impl std::fmt::Display for Key {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        let chars: String = self.0.iter().map(|&byte| byte as char).collect();
        write!(f, "{}", chars)
    }
}
impl std::fmt::Debug for Key {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        let chars: String = self.0.iter().map(|&byte| byte as char).collect();
        write!(f, "{}", chars)
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct KeyPair([u8; 4]);

impl KeyPair {
    fn new(a: Key, b: Key) -> Self {
        Self([a.0[0], a.0[1], b.0[0], b.0[1]])
    }
}

#[derive(Debug)]
struct Cave {
    valves: HashMap<Key, i32>,
    travel_cost: HashMap<KeyPair, i32>,
}

type In = Cave;
type Out = i32;
const PART1_RESULT: Out = 1651;
const PART2_RESULT: Out = 1707;

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut cave: Graph<Key, i32> = Graph::new();
    let mut valves: HashMap<Key, i32> = HashMap::new();
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
        .for_each(|v| {
            if v.data > 0 {
                valves.insert(v.key.clone(), v.data.clone());
            }
            cave.add_vertex(v);
        });
    let start = vec![Key::from("AA")];
    let travel_cost = valves
        .keys()
        .chain(start.iter())
        .flat_map(|a| {
            valves
                .keys()
                .map(|b| (KeyPair::new(*a, *b), cave.astar(a, b, |_| 1).unwrap().0))
        })
        .collect::<HashMap<_, _>>();

    Ok(Cave {
        valves,
        travel_cost,
    })
}

#[derive(Debug, Clone)]
struct State {
    current: Key,
    remaining: HashMap<Key, i32>, // (next, volume)
    pressure: i32,
    time: i32,
}

impl State {
    fn key(&self) -> String {
        let mut out = String::new();
        let mut r: Vec<String> = self.remaining.keys().map(|k| k.to_string()).collect();
        r.sort();
        format!(
            "{}:[{}]:{}:{}",
            self.current,
            r.join(":"),
            self.pressure,
            self.time
        )
    }
}

fn _max_possible(remaining: &HashSet<Key>, valves: Arc<HashMap<Key, i32>>) {
    let mut avail = valves
        .iter()
        .filter(|&(k, _)| remaining.contains(k))
        .collect::<Vec<_>>();
    println!("{:?}", avail);
}

fn search(
    s: State,
    time: i32,
    max: Arc<AtomicI32>,
    valves: Arc<HashMap<Key, i32>>,
    travel_cost: Arc<HashMap<KeyPair, i32>>,
    thread: bool,
) {
    if !s.remaining.is_empty() && s.time < time {
        let mut out: Vec<i32> = Vec::new();
        let mut handle = Vec::new();
        for (n, v) in s.remaining.iter() {
            let t_next = s.time + travel_cost[&KeyPair::new(s.current, *n)] + 1;
            let mut next = State {
                current: n.clone(),
                remaining: s.remaining.clone(),
                pressure: s.pressure + (time - t_next) * v,
                time: t_next,
            };
            next.remaining.remove(n);
            let valves = Arc::clone(&valves);
            let travel_cost = Arc::clone(&travel_cost);
            let max = Arc::clone(&max);
            if thread {
                // Only run top level search in threads
                handle.push(thread::spawn(move || {
                    search(next, time, max, valves, travel_cost, false);
                }));
            } else {
                search(next, time, max, valves, travel_cost, false);
            }
        }
        handle.into_iter().for_each(|h| h.join().unwrap());
    } else {
        max.fetch_max(s.pressure, Ordering::Relaxed);
    }
}

fn combinations<T: Clone>(v: &Vec<T>, n: usize) -> Vec<Vec<T>> {
    if n == 0 {
        return vec![vec![]];
    }

    if n > v.len() {
        return vec![];
    }

    let mut result = vec![];
    let mut stack: Vec<(usize, usize, Vec<T>)> = vec![(0, 0, vec![])];

    while let Some((i, count, mut current)) = stack.pop() {
        if count == n {
            result.push(current);
        } else {
            for j in i..v.len() {
                current.push(v[j].clone());
                stack.push((j + 1, count + 1, current.clone()));
                current.pop();
            }
        }
    }
    result
}

fn part1(
    Cave {
        valves,
        travel_cost,
    }: &In,
) -> Out {
    let start = State {
        current: Key::from("AA"),
        remaining: valves.clone(),
        pressure: 0,
        time: 0,
    };
    println!("Key:: {:?}", start.key());
    let valves = Arc::new(valves.clone());
    let travel_cost = Arc::new(travel_cost.clone());
    let max = Arc::new(AtomicI32::new(0));
    search(
        start,
        30,
        Arc::clone(&max),
        valves,
        travel_cost,
        std::env::var("THREAD").is_ok(),
    );
    (*max).load(Ordering::Relaxed)
}

fn part2(
    Cave {
        valves,
        travel_cost,
    }: &In,
) -> Out {
    let travel_cost = Arc::new(travel_cost.clone());
    let valves_arc = Arc::new(valves.clone());
    let v = valves.keys().cloned().collect::<Vec<Key>>();
    let mut max_both = 0;

    for i in ((valves.len() / 2) - 1)..=(valves.len() / 2) {
        for c in combinations(&v, i) {
            let mut v1: HashMap<Key, i32> = HashMap::new();
            let mut v2 = valves.clone();
            for k in c {
                v2.remove(&k);
                v1.insert(k, valves[&k]);
            }

            let s1 = v1.keys().cloned().collect::<Vec<_>>();
            let s2 = v2.keys().cloned().collect::<Vec<_>>();

            let start1 = State {
                current: Key::from("AA"),
                remaining: v1.clone(),
                pressure: 0,
                time: 0,
            };
            let start2 = State {
                current: Key::from("AA"),
                remaining: v2.clone(),
                pressure: 0,
                time: 0,
            };

            let max1 = Arc::new(AtomicI32::new(0));
            let max2 = Arc::new(AtomicI32::new(0));

            search(
                start1,
                26,
                Arc::clone(&max1),
                Arc::clone(&valves_arc),
                Arc::clone(&travel_cost),
                std::env::var("THREAD").is_ok(),
            );
            search(
                start2,
                26,
                Arc::clone(&max2),
                Arc::clone(&valves_arc),
                Arc::clone(&travel_cost),
                std::env::var("THREAD").is_ok(),
            );

            let max1 = (*max1).load(Ordering::Relaxed);
            let max2 = (*max2).load(Ordering::Relaxed);
            if max1 + max2 > max_both {
                println!("{:?} :: {:?}", s1, s2);
                println!(">> {} {} : {} -> {}\n", max1, max2, max_both, max1 + max2);
                max_both = max1 + max2;
            }
        }
    }
    max_both
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
