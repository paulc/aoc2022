use crate::Graph;
use crate::Vertex;

use std::collections::HashSet;
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
    I: Clone + Eq + Hash,
{
    type Item = I;
    fn next(&mut self) -> Option<Self::Item> {
        while let Some(i) = self.stack.pop() {
            if !self.discovered.contains(&i) {
                self.discovered.insert(i.clone());
                self.graph.get(&i).and_then(|v| {
                    for (e, _) in &v.edges {
                        self.stack.push(e.clone());
                    }
                    Some(())
                });
                return Some(i);
            }
        }
        None
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

    pub fn dfs<F>(&self, start: I, f: &mut F)
    where
        F: FnMut(&Vertex<I, D>),
    {
        let mut discovered: HashSet<I> = HashSet::new();
        Self::dfs_r(&self, &mut discovered, start, f);
    }

    fn dfs_r<F>(graph: &Graph<I, D>, discovered: &mut HashSet<I>, i: I, f: &mut F)
    where
        F: FnMut(&Vertex<I, D>),
    {
        discovered.insert(i.clone());
        graph.get(&i).and_then(|v| {
            f(v);
            for (e, _) in &v.edges {
                if !discovered.contains(e) {
                    Self::dfs_r(graph, discovered, e.clone(), f)
                }
            }
            Some(())
        });
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
        let mut out: Vec<String> = vec![];
        let mut f = |v: &Vertex<&str, &str>| out.push(v.key.to_string());
        g.dfs("A", &mut f);
        assert_eq!(
            out,
            vec!["A", "B", "D", "F", "E", "C", "G"]
                .iter()
                .map(|s| s.to_string())
                .collect::<Vec<_>>()
        );
    }
}
