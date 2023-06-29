use crate::Graph;

use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Debug;
use std::hash::Hash;

pub fn shortest_path<T>(graph: &Graph<T>, source: T, target: T) -> Option<(f64, Vec<T>)>
where
    T: Eq + PartialEq + Hash + Clone + Debug,
{
    let mut dist: HashMap<T, f64> = HashMap::new();
    let mut prev: HashMap<T, Option<T>> = HashMap::new();
    let mut q: HashSet<T> = HashSet::new();
    for v in graph.vertices() {
        dist.insert(v.clone(), f64::INFINITY);
        prev.insert(v.clone(), None);
        q.insert(v.clone());
    }
    dist.insert(source.clone(), 0.0);
    while !q.is_empty() {
        let u = dist
            .iter()
            .filter(|(k, _)| q.contains(k))
            .min_by(|a, b| a.1.partial_cmp(b.1).unwrap())
            .map(|(k, _)| k)
            .unwrap()
            .clone();

        if u == target {
            break;
        }

        q.remove(&u);
        for v in graph.neighbours(&u).unwrap() {
            if q.contains(&v.to) {
                let alt = dist.get(&u).unwrap() + v.cost;
                if alt < *dist.get(&v.to).unwrap() {
                    dist.insert(v.to.clone(), alt);
                    prev.insert(v.to.clone(), Some(u.clone()));
                }
            }
        }
    }
    if prev.contains_key(&target) {
        let mut path: Vec<T> = Vec::new();
        let mut u = target.clone();
        while let Some(v) = prev.get(&u).unwrap() {
            path.push(v.clone());
            u = v.clone();
        }
        Some((dist.get(&target).unwrap().clone(), path))
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
            assert_eq!(
                shortest_path(&g, start.to_string(), target.to_string()),
                Some((
                    path.0,
                    path.1
                        .to_vec()
                        .iter()
                        .map(|s| s.to_string())
                        .collect::<Vec<_>>()
                ))
            )
        }
    }
}
