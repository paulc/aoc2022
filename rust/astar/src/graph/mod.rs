pub mod dijkstra;

use std::collections::HashMap;
use std::hash::Hash;

// =========== Edge ===========

#[derive(Debug)]
pub struct Edge<T> {
    to: T,
    cost: f64,
}

impl<T: std::fmt::Display> std::fmt::Display for Edge<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}({:.1})", self.to, self.cost)
    }
}

// =========== Graph ===========

#[derive(Debug)]
pub struct Graph<T>(HashMap<T, Vec<Edge<T>>>);

impl<T: Eq + PartialEq + Hash + Clone> Graph<T> {
    pub fn new() -> Self {
        Self(HashMap::new())
    }
    pub fn add_vertex(&mut self, vertex: T) {
        self.0.entry(vertex).or_insert(Vec::new());
    }
    // We automatically add vertices referenced
    pub fn add_edge(&mut self, from: T, to: T, cost: f64) {
        self.0.entry(to.clone()).or_insert(Vec::new());
        self.0
            .entry(from)
            .or_insert(Vec::new())
            .push(Edge { to, cost });
    }
    pub fn add_edge_symmetric(&mut self, from: T, to: T, cost: f64) {
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

impl<T: std::fmt::Display> Graph<T> {
    fn to_dot(&self, name: String) -> String {
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

impl<T: Ord + std::fmt::Display> std::fmt::Display for Graph<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
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
mod tests {
    use super::*;

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
    fn test_graph() {
        assert_eq!(
            make_graph().to_string(),
            "AAA -> BBB(1.0), CCC(1.0), DDD(3.0)\nBBB -> DDD(2.0)\nCCC -> DDD(2.0)\nDDD -> AAA(3.0), EEE(4.0)\nEEE -> DDD(4.0)\n"
        )
    }
}
