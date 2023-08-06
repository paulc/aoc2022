use crate::Graph;

use std::collections::HashSet;
use std::hash::Hash;

pub struct DfsIter<'a, I, D>
where
    I: Clone + Copy + Eq + Hash,
{
    graph: &'a Graph<I, D>,
    discovered: HashSet<I>,
    stack: Vec<I>,
}

impl<'a, I, D> Iterator for DfsIter<'a, I, D>
where
    I: Clone + Copy + Eq + Hash,
{
    type Item = I;
    fn next(&mut self) -> Option<Self::Item> {
        while let Some(v) = self.stack.pop() {
            if !self.discovered.contains(&v) {
                self.discovered.insert(v);
                for (e, _) in &self.graph.get(&v).unwrap().edges {
                    self.stack.push(*e);
                }
                return Some(v);
            }
        }
        None
    }
}

impl<I, D> Graph<I, D>
where
    I: Clone + Copy + Eq + Hash,
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

    fn make_graph<'a>() -> Graph<&'a str, &'a str> {
        Graph::new_from_bidirectional_edges(vec![
            ("A", "B", 1),
            ("A", "C", 1),
            ("A", "E", 1),
            ("B", "D", 1),
            ("B", "F", 1),
            ("C", "G", 1),
            ("E", "F", 1),
        ])
    }

    #[test]
    fn test_dfs() {
        let g = make_graph();
        assert_eq!(
            g.dfs_iter("A")
                .map(|i| g.get(&i).and_then(|v| Some(v.key)))
                .collect::<Option<Vec<_>>>(),
            Some(vec!["A", "E", "F", "B", "D", "C", "G"])
        );
    }
}
