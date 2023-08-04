pub mod astar;
pub mod dfs;

use std::collections::HashMap;
use std::fmt::Debug;
use std::fmt::Display;
use std::hash::Hash;

#[derive(Debug, PartialEq)]
pub struct Vertex<I, D>(I, Option<D>, Vec<(I, i32)>)
where
    I: Clone + Eq + Hash;

impl<I, D> Vertex<I, D>
where
    I: Clone + Eq + Hash,
{
    pub fn new(id: I, data: Option<D>, edges: Vec<(I, i32)>) -> Self {
        Self(id, data, edges)
    }
    pub fn add_edge(&mut self, to: I, cost: i32) {
        self.2.push((to, cost))
    }
    pub fn key(&self) -> I {
        self.0.clone()
    }
    pub fn data(&self) -> Option<&D> {
        self.1.as_ref()
    }
    pub fn data_mut(&mut self) -> Option<&mut D> {
        self.1.as_mut()
    }
    pub fn edges(&self) -> impl Iterator<Item = (I, i32)> + '_ {
        self.2.iter().cloned()
    }
    pub fn edges_mut(&mut self) -> &mut Vec<(I, i32)> {
        &mut self.2
    }
}

impl<I, D> Display for Vertex<I, D>
where
    I: Display + Clone + Eq + Hash,
    D: Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{} <{}> -> ",
            self.0,
            match &self.1 {
                Some(d) => d.to_string(),
                None => "".to_string(),
            }
        )?;
        let mut n = self.2.len();
        for (d, c) in &self.2 {
            write!(f, "[{}]({})", d, c)?;
            n -= 1;
            if n > 0 {
                write!(f, ",")?;
            }
        }
        Ok(())
    }
}

#[derive(Debug, PartialEq)]
pub struct Graph<I, D>(HashMap<I, Vertex<I, D>>)
where
    I: Clone + Eq + Hash;

impl<I, D> Graph<I, D>
where
    I: Clone + Eq + Hash,
{
    pub fn new() -> Self {
        Self(HashMap::new())
    }
    pub fn new_from_edges(edges: Vec<(I, I, i32)>) -> Self {
        let mut out = Self::new();
        for (v1, v2, cost) in edges {
            out.0
                .entry(v2.clone())
                .or_insert_with(|| Vertex::new(v2.clone(), None, vec![]));
            out.0
                .entry(v1.clone())
                .or_insert_with(|| Vertex::new(v1.clone(), None, vec![]))
                .add_edge(v2, cost);
        }
        out
    }
    pub fn add_vertex(&mut self, v: Vertex<I, D>) {
        self.0.entry(v.0.clone()).or_insert_with(|| v);
    }
    pub fn vertices(&self) -> impl Iterator<Item = &Vertex<I, D>> {
        self.0.values()
    }
    pub fn get(&self, key: &I) -> Option<&Vertex<I, D>> {
        self.0.get(key)
    }
    pub fn get_mut(&mut self, key: &I) -> Option<&mut Vertex<I, D>> {
        self.0.get_mut(key)
    }
}

impl<I, D> Display for Graph<I, D>
where
    I: Display + Clone + Eq + Hash,
    D: Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for v in self.vertices() {
            writeln!(f, "{}", v)?;
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[derive(Debug, Clone, PartialEq)]
    struct E(());

    impl Display for E {
        fn fmt(&self, _f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            Ok(())
        }
    }

    #[derive(Debug, PartialEq)]
    struct D(i32, i32);
    impl Display for D {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "({},{})", self.0, self.1)
        }
    }

    fn make_graph() -> Graph<&'static str, E> {
        Graph::new_from_edges(vec![
            ("DD", "EE", 1),
            ("BB", "DD", 5),
            ("AA", "BB", 1),
            ("AA", "CC", 2),
            ("CC", "DD", 1),
        ])
    }

    #[test]
    fn test_graph_new() {
        let mut g: Graph<&'static str, E> = Graph::new();
        g.add_vertex(Vertex::new("AA", None, vec![("BB", 1), ("CC", 2)]));
        g.add_vertex(Vertex::new("BB", None, vec![("DD", 5)]));
        g.add_vertex(Vertex::new("CC", None, vec![("DD", 1)]));
        g.add_vertex(Vertex::new("DD", None, vec![("EE", 1)]));
        g.add_vertex(Vertex::new("EE", None, vec![]));
        assert_eq!(g, make_graph());
    }

    #[test]
    fn test_graph_from_edges() {
        let g = make_graph();
        let s = g.to_string();
        assert_eq!(
            {
                let mut l = s.lines().collect::<Vec<_>>();
                l.sort();
                l
            },
            vec![
                "AA <> -> [BB](1),[CC](2)",
                "BB <> -> [DD](5)",
                "CC <> -> [DD](1)",
                "DD <> -> [EE](1)",
                "EE <> -> "
            ]
        );
    }

    #[test]
    fn test_graph_get() {
        let g = make_graph();
        assert_eq!(
            g.get(&"AA"),
            Some(&Vertex::new("AA", None, vec![("BB", 1), ("CC", 2)]))
        );
        assert_eq!(g.get(&"ZZ"), None);
    }

    #[test]
    fn test_graph_vertices() {
        let g = make_graph();
        assert_eq!(
            {
                let mut v = g.vertices().map(|v| v.key()).collect::<Vec<_>>();
                v.sort();
                v
            },
            vec!["AA", "BB", "CC", "DD", "EE"]
        );
    }

    #[test]
    fn test_graph_add_vertex() {
        let mut g = make_graph();
        assert_eq!(g.get(&"ZZ"), None);
        g.add_vertex(Vertex::new("ZZ", None, vec![("AA", 99)]));
        assert_eq!(
            {
                let mut v = g.vertices().map(|v| v.key()).collect::<Vec<_>>();
                v.sort();
                v
            },
            vec!["AA", "BB", "CC", "DD", "EE", "ZZ"]
        );
        assert_eq!(
            g.get(&"ZZ"),
            Some(&Vertex::new("ZZ", None, vec![("AA", 99)]))
        );
    }

    #[test]
    fn test_graph_get_mut() {
        let mut g = make_graph();
        g.get_mut(&"AA").unwrap().add_edge("EE", 10);
        assert_eq!(
            g.get(&"AA"),
            Some(&Vertex::new(
                "AA",
                None,
                vec![("BB", 1), ("CC", 2), ("EE", 10)]
            ))
        );
    }

    #[test]
    fn test_graph_get_mut2() {
        let mut g: Graph<&'static str, i32> = Graph::new();
        g.add_vertex(Vertex::new("AA", Some(0), vec![]));
        g.get_mut(&"AA").unwrap().1 = Some(99);
        assert_eq!(g.get(&"AA").unwrap().data(), Some(&99));
    }

    #[test]
    fn test_vertex_add_edge() {
        let mut g = make_graph();
        g.get_mut(&"AA").unwrap().add_edge("ZZ", 99);
        assert_eq!(
            g.get(&"AA").unwrap().edges().collect::<Vec<_>>(),
            vec![("BB", 1), ("CC", 2), ("ZZ", 99)]
        );
    }

    #[test]
    fn test_vertex_key() {
        let g = make_graph();
        assert_eq!(g.get(&"AA").unwrap().key(), "AA");
    }

    #[test]
    fn test_vertex_data() {
        let mut g: Graph<&'static str, D> = Graph::new();
        g.add_vertex(Vertex::new("AA", Some(D(0, 1)), vec![]));
        assert_eq!(g.get(&"AA").unwrap().data(), Some(&D(0, 1)));
    }

    #[test]
    fn test_vertex_data_mut() {
        let mut g: Graph<&'static str, D> = Graph::new();
        g.add_vertex(Vertex::new("AA", Some(D(0, 1)), vec![]));
        g.get_mut(&"AA").unwrap().data_mut().unwrap().0 = 99;
        assert_eq!(g.get(&"AA").unwrap().data(), Some(&D(99, 1)));
    }

    #[test]
    fn test_vertex_edges() {
        let g = make_graph();
        assert_eq!(
            g.get(&"AA").unwrap().edges().collect::<Vec<_>>(),
            vec![("BB", 1), ("CC", 2)]
        );
    }

    #[test]
    fn test_vertex_edges_mut() {
        let mut g = make_graph();
        g.get_mut(&"AA").unwrap().edges_mut().push(("ZZ", 99));
        assert_eq!(
            g.get(&"AA").unwrap().edges().collect::<Vec<_>>(),
            vec![("BB", 1), ("CC", 2), ("ZZ", 99)]
        );
    }
}
