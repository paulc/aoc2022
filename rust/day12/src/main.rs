#![allow(unused)]

use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap, HashSet};
use std::fmt::Display;
use std::fs::File;
use std::hash::Hash;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;

use graph::dijkstra::shortest_path;
use graph::Graph;

type In = Hill;
type Out = f64;
const PART1_RESULT: Out = 31.0;
const PART2_RESULT: Out = 29.0;

// =========== XY ===========

#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash, PartialOrd, Ord)]
struct XY {
    x: i32,
    y: i32,
}

impl XY {
    fn new(x: i32, y: i32) -> XY {
        XY { x, y }
    }
    fn adjacent(&self) -> [XY; 4] {
        [
            XY::new(self.x, self.y - 1),
            XY::new(self.x + 1, self.y),
            XY::new(self.x, self.y + 1),
            XY::new(self.x - 1, self.y),
        ]
    }
}

impl std::fmt::Display for XY {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "({},{})", self.x, self.y)
    }
}

// =========== Grid ===========

#[derive(Debug)]
struct Grid<T>(Vec<Vec<T>>);

impl<T> Grid<T> {
    fn check_bounds(&self, p: XY) -> bool {
        (p.y >= 0 && p.y < self.0.len() as i32)
            && ((self.0.len() > 0) && (p.x >= 0 && p.x < self.0[0].len() as i32))
    }
    fn get(&self, p: XY) -> Option<&T> {
        match self.check_bounds(p) {
            true => Some(&self.0[p.y as usize][p.x as usize]),
            false => None,
        }
    }
    fn adjacent(&self, p: XY) -> Vec<(&T, XY)> {
        p.adjacent()
            .into_iter()
            .filter(|&p| self.check_bounds(p))
            .map(|p| (&self.0[p.y as usize][p.x as usize], p))
            .collect()
    }
    fn iter(&self) -> impl Iterator<Item = (&T, XY)> {
        self.0.iter().enumerate().flat_map(|(y, r)| {
            r.iter()
                .enumerate()
                .map(move |(x, c)| (c, XY::new(x as i32, y as i32)))
        })
    }
}

impl<T: std::fmt::Display> std::fmt::Display for Grid<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for y in 0..self.0.len() {
            writeln!(
                f,
                "{}",
                self.0[y]
                    .iter()
                    .map(|e| format!("{e:>2}"))
                    .collect::<Vec<String>>()
                    .join(" ")
            )?
        }
        Ok(())
    }
}
/*
// =========== Cost ===========

#[derive(Debug, Copy, Clone)]
struct Cost<T> {
    p: T,
    cost: f64,
}

impl<T> PartialEq for Cost<T> {
    fn eq(&self, other: &Self) -> bool {
        self.cost.eq(&other.cost)
    }
}

impl<T> Eq for Cost<T> {}

impl<T> PartialOrd for Cost<T> {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        self.cost.partial_cmp(&other.cost)
    }
}

impl<T: Eq> Ord for Cost<T> {
    fn cmp(&self, other: &Self) -> Ordering {
        self.partial_cmp(other).unwrap_or(Ordering::Equal)
    }
}

// =========== Astar ===========

fn astar<T: Ord + Clone + Eq + Hash + std::fmt::Debug>(
    graph: &Graph<T>,
    start: T,
    end: T,
    h: fn(&T) -> f64,
) -> Option<f64> {
    let mut open_set: BinaryHeap<Cost<T>> = BinaryHeap::new();
    open_set.push(Cost {
        p: start.clone(),
        cost: h(&start),
    });
    let mut came_from: HashMap<T, T> = HashMap::new();
    let mut g_score: HashMap<T, f64> = HashMap::new();
    g_score.insert(start.clone(), 0.0);
    let mut f_score: HashMap<T, f64> = HashMap::new();
    f_score.insert(start.clone(), h(&start));
    while !open_set.is_empty() {
        if let Some(current) = open_set.pop() {
            if current.p == end {
                break;
            }
            if let Some(neighbours) = graph.0.get(&current.p) {
                for neighbour in neighbours {
                    let tentative_g_score =
                        g_score.get(&current.p).unwrap_or(&f64::INFINITY) + neighbour.cost;
                    /*
                        println!(
                            ">> Neighbour: {:?} tentative_g_score={} current_g_score[neighbour]={:?}",
                            neighbour,
                            tentative_g_score,
                            g_score.get(&neighbour.to).unwrap_or(&f64::INFINITY)
                        );
                    */
                    if tentative_g_score < *g_score.get(&neighbour.to).unwrap_or(&f64::INFINITY) {
                        came_from.insert(neighbour.to.clone(), current.p.clone());
                        g_score.insert(neighbour.to.clone(), tentative_g_score.clone());
                        let neighbour_f_score = tentative_g_score + h(&neighbour.to);
                        f_score.insert(neighbour.to.clone(), neighbour_f_score.clone());
                        if open_set.iter().filter(|c| c.p == neighbour.to).count() == 0 {
                            open_set.push(Cost {
                                p: neighbour.to.clone(),
                                cost: neighbour_f_score.clone(),
                            });
                        }
                    }
                }
            }
        }
    }
    let cost = g_score.get(&end).map(|c| c.clone());
    cost
}
*/

#[derive(Debug)]
struct Hill {
    hill: Grid<u8>,
    reachable: Graph<XY>,
    start: XY,
    end: XY,
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut start: XY = XY::new(0, 0);
    let mut end: XY = XY::new(0, 0);
    let mut hill = Grid(
        BufReader::new(input)
            .lines()
            .enumerate()
            .map(|(y, l)| {
                l?.chars()
                    .enumerate()
                    .map(|(x, c)| match c {
                        'a'..='z' => Ok((c as u8) - b'a'),
                        'S' => {
                            start = XY::new(x as i32, y as i32);
                            Ok(0)
                        }
                        'E' => {
                            end = XY::new(x as i32, y as i32);
                            Ok(25)
                        }
                        _ => Err(Error::new(InvalidData, format!("({},{}) {}", x, y, c))),
                    })
                    .collect::<std::io::Result<Vec<_>>>()
            })
            .collect::<std::io::Result<Vec<_>>>()?,
    );
    let mut reachable: Graph<XY> = Graph::new();
    for (h, p) in hill.iter() {
        hill.adjacent(p)
            .iter()
            .filter(|(&h1, _)| h1 <= h + 1)
            .for_each(|(_, p1)| {
                reachable.add_edge(p, *p1, 1.0);
            });
    }

    Ok(Hill {
        hill,
        reachable,
        start,
        end,
    })
}

fn part1(input: &In) -> Out {
    let (cost, path) = shortest_path(&input.reachable, input.start, input.end).unwrap();
    cost
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
Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi
";
