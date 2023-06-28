use crate::graph::Graph;

use std::collections::HashMap;
use std::collections::HashSet;
use std::fmt::Debug;
use std::hash::Hash;

pub fn dijkstra<T>(graph: Graph<T>, source: T, target: T) -> Option<(f64, Vec<T>)>
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
                println!(">> {:?}", v);
                let alt = dist.get(&u).unwrap() + v.cost;
                if alt < *dist.get(&v.to).unwrap() {
                    dist.insert(v.to.clone(), alt);
                    prev.insert(v.to.clone(), Some(u.clone()));
                }
            }
        }
    }
    println!("Dist: {:?}", dist);
    println!("Prev: {:?}", prev);
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
    use crate::graph::Graph;

    fn make_graph() -> Graph<String> {
        let mut g: Graph<String> = Graph::new();
        g.add_vertex("AAA".to_string());
        g.add_vertex("BBB".to_string());
        g.add_edge("AAA".to_string(), "BBB".to_string(), 1.0);
        g.add_edge("AAA".to_string(), "CCC".to_string(), 1.0);
        g.add_edge("BBB".to_string(), "DDD".to_string(), 2.0);
        g.add_edge("CCC".to_string(), "DDD".to_string(), 2.0);
        g.add_edge_symmetric("DDD".to_string(), "AAA".to_string(), 3.0);
        g.add_edge_symmetric("DDD".to_string(), "EEE".to_string(), 4.0);
        g
    }

    #[test]
    fn test_dijkstra() {
        let g = make_graph();
        println!("{}", g);
        dbg!(dijkstra(g, "EEE".to_string(), "CCC".to_string()));
    }
}
