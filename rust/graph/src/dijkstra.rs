use crate::Graph;

use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::collections::HashMap;
use std::fmt::Debug;
use std::hash::Hash;

#[derive(Debug, Clone)]
struct State<'a, T> {
    element: &'a T,
    cost: f64,
}

impl<'a, T> State<'a, T> {
    fn new(element: &'a T, cost: f64) -> Self {
        State { element, cost }
    }
}

impl<'a, T> PartialEq for State<'a, T> {
    fn eq(&self, other: &Self) -> bool {
        self.cost == other.cost
    }
}

impl<'a, T> Eq for State<'a, T> {}

impl<'a, T> PartialOrd for State<'a, T> {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}
impl<'a, T> Ord for State<'a, T> {
    // We need to extract lowest cost item from priq so invert ordering
    fn cmp(&self, other: &Self) -> Ordering {
        other.cost.partial_cmp(&self.cost).unwrap()
    }
}

pub fn shortest_path_all<'a, T>(graph: &'a Graph<T>, source: &'a T) -> Vec<(&'a T, f64)>
where
    T: Eq + Hash + Clone,
{
    let mut q: BinaryHeap<State<T>> = BinaryHeap::new();
    let mut dist: HashMap<&T, f64> = HashMap::new();
    dist.insert(&source, 0.0);
    q.push(State::new(&source, 0.0));

    while let Some(u) = q.pop() {
        if let Some(edges) = graph.neighbours(u.element) {
            for edge in edges {
                let alt = dist.get(u.element).unwrap_or(&f64::INFINITY) + edge.cost;
                if alt < *dist.get(&edge.to).unwrap_or(&f64::INFINITY) {
                    dist.insert(&edge.to, alt);
                    q.push(State::new(&edge.to, alt))
                }
            }
        }
    }
    dist.into_iter().collect::<Vec<_>>()
}

pub fn shortest_path<T>(graph: &Graph<T>, source: T, target: T) -> Option<(f64, Vec<T>)>
where
    T: Eq + Hash + Clone,
{
    let mut q: BinaryHeap<State<T>> = BinaryHeap::new();
    let mut dist: HashMap<&T, f64> = HashMap::new();
    let mut prev: HashMap<&T, Option<&T>> = HashMap::new();
    dist.insert(&source, 0.0);
    q.push(State::new(&source, 0.0));

    while let Some(u) = q.pop() {
        if *u.element == target {
            break;
        }
        if let Some(edges) = graph.neighbours(u.element) {
            for edge in edges {
                let alt = dist.get(u.element).unwrap_or(&f64::INFINITY) + edge.cost;
                if alt < *dist.get(&edge.to).unwrap_or(&f64::INFINITY) {
                    dist.insert(&edge.to, alt);
                    prev.insert(&edge.to, Some(u.element));
                    q.push(State::new(&edge.to, alt))
                }
            }
        }
    }
    if let Some(cost) = dist.get(&target) {
        let mut path: Vec<T> = Vec::new();
        let mut u = target;
        while let Some(Some(v)) = prev.get(&u) {
            let v = v.clone();
            path.push(v.clone());
            u = v.clone();
        }
        Some((*cost, path))
    } else {
        None
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::*;

    #[test]
    fn test_shortest_path() {
        let g = make_graph();
        for (start, target, path) in [
            ("EEE", "CCC", (8.0, ["AAA", "DDD", "EEE"])),
            ("BBB", "CCC", (6.0, ["AAA", "DDD", "BBB"])),
            ("AAA", "EEE", (6.0, ["DDD", "CCC", "AAA"])),
        ] {
            let path_vec = path
                .1
                .to_vec()
                .iter()
                .map(|s| s.to_string())
                .collect::<Vec<_>>();
            assert_eq!(
                shortest_path(&g, start.to_string(), target.to_string()),
                Some((path.0, path_vec))
            )
        }
    }

    #[test]
    fn test_shortest_path_all() {
        let g = make_graph();
        let source = "CCC".to_string();
        let mut all = shortest_path_all(&g, &source);
        all.sort_by_key(|n| n.0);
        assert_eq!(
            all,
            vec![
                (&"AAA".to_string(), 4.0),
                (&"BBB".to_string(), 5.0),
                (&"CCC".to_string(), 0.0),
                (&"DDD".to_string(), 1.0),
                (&"EEE".to_string(), 5.0)
            ]
        );
    }

    #[test]
    fn test_state_order() {
        let (e1, e2) = (String::from("AAA"), String::from("BBB"));
        let (s1, s2) = (State::new(&e1, 1.0), State::new(&e2, 2.0));
        assert_eq!(s1.cmp(&s2), Ordering::Greater);
    }

    #[test]
    fn test_priq() {
        let (e1, e2, e3, e4) = (
            String::from("AAA"),
            String::from("BBB"),
            String::from("CCC"),
            String::from("DDD"),
        );
        let (s1, s2, s3, s4) = (
            State::new(&e1, 3.0),
            State::new(&e2, 10.0),
            State::new(&e3, 1.0),
            State::new(&e4, 5.0),
        );
        let mut q: BinaryHeap<State<String>> = BinaryHeap::new();
        q.push(s1);
        q.push(s2);
        q.push(s3);
        assert_eq!(q.pop().unwrap().element, "CCC");
        q.push(s4);
        assert_eq!(q.pop().unwrap().element, "AAA");
        assert_eq!(q.pop().unwrap().element, "DDD");
        assert_eq!(q.pop().unwrap().element, "BBB");
        assert!(q.is_empty());
    }
}
