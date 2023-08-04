use crate::{Graph, Vertex};

use std::collections::HashSet;
// use std::collections::VecDeque;
use std::hash::Hash;

pub struct DfsIter<'a, I, D>
where
    I: Clone + Eq + Hash,
{
    graph: &'a Graph<I, D>,
    discovered: HashSet<I>,
    stack: Vec<I>,
}

impl<'a, I, D> Iterator for DfsIter<'a, I, D>
where
    I: Clone + Eq + Hash + std::fmt::Debug,
    D: std::fmt::Debug,
{
    type Item = &'a Vertex<I, D>;
    fn next(&mut self) -> Option<Self::Item> {
        if let Some(v) = self.stack.pop().and_then(|v| self.graph.get(&v)) {
            println!("{:?}", v);
            if !self.discovered.contains(&v.0) {
                self.discovered.insert(v.0.clone());
                for (e, _) in v.edges() {
                    self.stack.push(e.clone());
                }
            }
            Some(v)
        } else {
            None
        }
    }
}

impl<I, D> Graph<I, D>
where
    I: Clone + Eq + Hash,
{
    pub fn dfs_iter(&self, root: I) -> DfsIter<I, D> {
        DfsIter {
            graph: &self,
            discovered: HashSet::new(),
            stack: vec![root],
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::*;

    #[test]
    fn test_astar_simple() {
        let _g: Graph<&str, ()> = Graph::new_from_edges(vec![
            ("A", "B", 1),
            ("A", "C", 1),
            ("A", "E", 1),
            ("B", "D", 1),
            ("B", "F", 1),
            ("C", "G", 1),
            ("E", "F", 1),
        ]);
    }
}
