pub mod dijkstra;

use std::collections::HashMap;
use std::fmt;
use std::fmt::Display;
use std::fmt::Formatter;
use std::hash::Hash;

// =========== Edge ===========

#[derive(Debug)]
pub struct Edge<T> {
    to: T,
    cost: f64,
}

impl<T: Display> Display for Edge<T> {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        write!(f, "{}({:.1})", self.to, self.cost)
    }
}

// =========== Graph ===========

#[derive(Debug)]
pub struct Graph<T>(HashMap<T, Vec<Edge<T>>>);

impl<T> Graph<T>
where
    T: Eq + PartialEq + Hash + Clone,
{
    pub fn new() -> Self {
        Self(HashMap::new())
    }
    pub fn add_vertex(&mut self, vertex: T) {
        self.0.entry(vertex).or_insert(Vec::new());
    }
    pub fn add_edge(&mut self, from: T, to: T, cost: f64) {
        // We automatically add vertices referenced
        self.0.entry(to.clone()).or_insert(Vec::new());
        self.0
            .entry(from)
            .or_insert(Vec::new())
            .push(Edge { to, cost });
    }
    pub fn add_edge_symmetric(&mut self, from: T, to: T, cost: f64) {
        // We automatically add vertices referenced
        self.0.entry(from.clone()).or_insert(Vec::new()).push(Edge {
            to: to.clone(),
            cost: cost.clone(),
        });
        self.0.entry(to.clone()).or_insert(Vec::new()).push(Edge {
            to: from.clone(),
            cost: cost.clone(),
        });
    }
    pub fn vertices(&self) -> impl Iterator<Item = &T> {
        self.0.iter().map(|(k, _)| k)
    }
    pub fn neighbours(&self, target: &T) -> Option<&Vec<Edge<T>>> {
        self.0.get(target)
    }
}

impl<T: Display> Graph<T> {
    pub fn to_dot(&self, name: String) -> String {
        let mut s: Vec<String> = Vec::new();
        s.push(format!("digraph {} {{", name));
        for (vertex, edges) in self.0.iter() {
            for edge in edges {
                s.push(format!("{} -> {} [label={}]", vertex, edge.to, edge.cost));
            }
        }
        s.push("}\n".to_string());
        s.join("\n")
    }
}

impl<T: Ord + Display> Display for Graph<T> {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        let mut keys = self.0.iter().collect::<Vec<_>>();
        keys.sort_by_key(|(k, _)| k.clone());
        for (k, v) in keys {
            writeln!(
                f,
                "{} -> {}",
                k,
                v.iter()
                    .map(|e| e.to_string())
                    .collect::<Vec<_>>()
                    .join(", ")
            )?;
        }
        Ok(())
    }
}

#[cfg(test)]
fn make_graph() -> Graph<String> {
    let mut g: Graph<String> = Graph::new();
    g.add_vertex("AAA".to_string());
    g.add_vertex("BBB".to_string());
    g.add_edge("AAA".to_string(), "BBB".to_string(), 1.0);
    g.add_edge("AAA".to_string(), "CCC".to_string(), 1.0);
    g.add_edge("BBB".to_string(), "DDD".to_string(), 2.0);
    g.add_edge("CCC".to_string(), "DDD".to_string(), 1.0);
    g.add_edge_symmetric("DDD".to_string(), "AAA".to_string(), 3.0);
    g.add_edge_symmetric("DDD".to_string(), "EEE".to_string(), 4.0);
    g
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_to_string() {
        assert_eq!(
            make_graph().to_string(),
            "AAA -> BBB(1.0), CCC(1.0), DDD(3.0)\nBBB -> DDD(2.0)\nCCC -> DDD(1.0)\nDDD -> AAA(3.0), EEE(4.0)\nEEE -> DDD(4.0)\n"
        )
    }

    #[test]
    fn test_vertices() {
        let g = make_graph();
        let mut vertices: Vec<&String> = g.vertices().collect();
        vertices.sort();
        assert_eq!(vertices, vec!["AAA", "BBB", "CCC", "DDD", "EEE"])
    }

    #[test]
    fn test_neighbours() {
        let g = make_graph();
        for (v, expected) in [
            ("AAA", vec!["BBB", "CCC", "DDD"]),
            ("DDD", vec!["AAA", "EEE"]),
        ] {
            let mut neighbours = g
                .neighbours(&v.to_string())
                .unwrap()
                .iter()
                .map(|e| e.to.clone())
                .collect::<Vec<String>>();
            neighbours.sort();
            assert_eq!(neighbours, expected);
        }
    }

    #[test]
    fn test_dot() {
        println!("{}", make_graph().to_dot("test".to_string()))
    }
}
