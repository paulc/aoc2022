use crate::Graph;

use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap};
use std::hash::Hash;

#[derive(Debug, PartialEq, Eq)]
struct V<I>(I, i32)
where
    I: Eq;

impl<I> PartialOrd for V<I>
where
    I: Eq,
{
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        self.1.partial_cmp(&other.1).map(|o| match o {
            Ordering::Less => Ordering::Greater,
            Ordering::Equal => Ordering::Equal,
            Ordering::Greater => Ordering::Less,
        })
    }
}
impl<I> Ord for V<I>
where
    I: Eq,
{
    fn cmp(&self, other: &Self) -> Ordering {
        match self.1.cmp(&other.1) {
            Ordering::Less => Ordering::Greater,
            Ordering::Equal => Ordering::Equal,
            Ordering::Greater => Ordering::Less,
        }
    }
}

impl<I, D> Graph<I, D>
where
    I: Clone + Eq + Hash,
{
    pub fn astar<F>(&self, start: I, target: I, h: F) -> Option<(i32, Vec<I>)>
    where
        F: Fn(&I) -> i32,
    {
        let mut open: BinaryHeap<V<I>> = BinaryHeap::new();
        let mut from: HashMap<I, I> = HashMap::new();
        let mut score: HashMap<I, i32> = HashMap::new();
        open.push(V(start.clone(), h(&start)));
        score.insert(start.clone(), 0);
        while let Some(current) = open.pop() {
            if current.0 == target {
                if let Some(cost) = score.get(&target) {
                    let mut current = current.0;
                    let mut path = vec![current.clone()];
                    while let Some(prev) = from.get(&current) {
                        path.push(prev.clone());
                        current = prev.clone();
                    }
                    path.reverse();
                    return Some((*cost, path));
                } else {
                    return None;
                }
            }
            for (n, d) in self.get(&current.0).unwrap().edges() {
                let tentative = score[&current.0] + d;
                if tentative < *score.get(&n).unwrap_or(&i32::MAX) {
                    from.insert(n.clone(), current.0.clone());
                    score.insert(n.clone(), tentative);
                    open.push(V(n.clone(), tentative + h(&n)));
                }
            }
        }
        None
    }
}

#[cfg(test)]
mod tests {
    use crate::*;
    use std::cmp::{max, min};
    use std::fs;

    #[test]
    fn test_astar_simple() {
        let g: Graph<&str, ()> = Graph::new_from_edges(vec![
            ("A", "B", 2),
            ("A", "C", 3),
            ("B", "D", 10),
            ("B", "E", 3),
            ("C", "D", 3),
            ("C", "E", 5),
            ("D", "F", 1),
            ("E", "F", 1),
        ]);
        assert_eq!(
            g.astar("A", "F", |_| 1),
            Some((6, vec!["A", "B", "E", "F"]))
        );
    }

    // From aoc2021/day15
    fn make_graph(path: &str) -> Graph<(usize, usize), ()> {
        let a = fs::read_to_string(path)
            .unwrap()
            .lines()
            .map(|l| l.as_bytes().iter().map(|b| *b - b'0').collect::<Vec<_>>())
            .collect::<Vec<_>>();
        let mut g: Graph<(usize, usize), ()> = Graph::new();
        for y in 0..a.len() {
            for x in 0..a[0].len() {
                g.add_vertex(Vertex::new(
                    (x, y),
                    None,
                    adj(&a, (x, y))
                        .iter()
                        .map(|&(x, y)| ((x, y), a[y][x] as i32))
                        .collect::<Vec<_>>(),
                ))
            }
        }
        g
    }

    fn adj(a: &Vec<Vec<u8>>, (x, y): (usize, usize)) -> Vec<(usize, usize)> {
        let (x_max, y_max) = (a[0].len() as i32, a.len() as i32);
        let (x, y) = (x as i32, y as i32);
        let mut out = vec![];
        for (dx, dy) in vec![(-1, 0), (1, 0), (0, -1), (0, 1)] {
            if x + dx >= 0 && x + dx < x_max && y + dy >= 0 && y + dy < y_max {
                out.push(((x + dx) as usize, (y + dy) as usize));
            }
        }
        out
    }

    fn md((x1, y1): (usize, usize), (x2, y2): (usize, usize)) -> usize {
        (max(x1, x2) - min(x1, x2)) + (max(y1, y2) - min(y1, y2))
    }

    #[test]
    fn test_astar_grid() {
        let g = make_graph("testdata/grid.txt");
        let (cost, path) = g
            .astar((0, 0), (9, 9), |&(x, y)| md((x, y), (9, 9)) as i32)
            .unwrap();
        assert_eq!((cost, path.len()), (40, 19));
    }

    #[test]
    fn test_astar_grid_large() {
        let g = make_graph("testdata/grid_large.txt");
        let (cost, _path) = g
            .astar((0, 0), (99, 99), |&(x, y)| md((x, y), (99, 99)) as i32)
            .unwrap();
        assert_eq!(cost, 602);
    }
}
