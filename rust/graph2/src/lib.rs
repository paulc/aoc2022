pub mod astar;
pub mod dfs;

use std::collections::HashMap;
use std::fmt::Debug;
use std::fmt::Display;
use std::hash::Hash;

#[derive(Debug, PartialEq)]
pub struct Vertex<I, D>
where
    I: Clone + Copy + Eq + Hash,
{
    key: I,
    data: Option<D>,
    edges: Vec<(I, i32)>,
}

impl<I, D> Vertex<I, D>
where
    I: Clone + Copy + Eq + Hash,
{
    pub fn new(key: I, data: Option<D>, edges: Vec<(I, i32)>) -> Self {
        Self { key, data, edges }
    }
    pub fn add_edge(&mut self, to: I, cost: i32) {
        self.edges.push((to, cost))
    }
}

impl<I, D> Display for Vertex<I, D>
where
    I: Display + Clone + Copy + Eq + Hash,
    D: Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{} <{}> -> ",
            self.key,
            match &self.data {
                Some(d) => d.to_string(),
                None => "".to_string(),
            }
        )?;
        let mut n = self.edges.len();
        for (d, c) in &self.edges {
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
    I: Clone + Copy + Eq + Hash;

impl<I, D> Graph<I, D>
where
    I: Clone + Copy + Eq + Hash,
{
    pub fn new() -> Self {
        Self(HashMap::new())
    }
    pub fn new_from_edges(edges: Vec<(I, I, i32)>) -> Self {
        let mut out = Self::new();
        for (v1, v2, cost) in edges {
            out.0
                .entry(v2)
                .or_insert_with(|| Vertex::new(v2, None, vec![]));
            out.0
                .entry(v1)
                .or_insert_with(|| Vertex::new(v1, None, vec![]))
                .add_edge(v2, cost);
        }
        out
    }
    pub fn new_from_bidirectional_edges(edges: Vec<(I, I, i32)>) -> Self {
        let mut out = Self::new();
        for (v1, v2, cost) in edges {
            out.0
                .entry(v2)
                .or_insert_with(|| Vertex::new(v2, None, vec![]))
                .add_edge(v1, cost);
            out.0
                .entry(v1)
                .or_insert_with(|| Vertex::new(v1, None, vec![]))
                .add_edge(v2, cost);
        }
        out
    }
    pub fn add_vertex(&mut self, v: Vertex<I, D>) {
        self.0.entry(v.key).or_insert_with(|| v);
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
    I: Display + Clone + Copy + Eq + Hash,
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
                let mut v = g.vertices().map(|v| v.key).collect::<Vec<_>>();
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
                let mut v = g.vertices().map(|v| v.key).collect::<Vec<_>>();
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
        g.get_mut(&"AA").and_then(|v| {
            v.add_edge("EE", 10);
            Some(())
        });
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
        g.get_mut(&"AA").and_then(|v| {
            v.data = Some(99);
            Some(())
        });
        assert_eq!(g.get(&"AA").and_then(|v| v.data), Some(99));
    }

    #[test]
    fn test_vertex_add_edge() {
        let mut g = make_graph();
        g.get_mut(&"AA").and_then(|v| {
            v.add_edge("ZZ", 99);
            Some(())
        });
        assert_eq!(
            g.get(&"AA")
                .and_then(|v| Some(v.edges.iter().collect::<Vec<_>>())),
            Some(vec![&("BB", 1), &("CC", 2), &("ZZ", 99)])
        );
    }

    #[test]
    fn test_vertex_key() {
        let g = make_graph();
        assert_eq!(g.get(&"AA").and_then(|v| Some(v.key)), Some("AA"));
    }

    #[test]
    fn test_vertex_data() {
        let mut g: Graph<&'static str, D> = Graph::new();
        g.add_vertex(Vertex::new("AA", Some(D(0, 1)), vec![]));
        assert_eq!(g.get(&"AA").and_then(|v| v.data.as_ref()), Some(&D(0, 1)));
    }

    #[test]
    fn test_vertex_data_mut() {
        let mut g: Graph<&'static str, D> = Graph::new();
        g.add_vertex(Vertex::new("AA", Some(D(0, 1)), vec![]));
        if let Some(d) = g.get_mut(&"AA").and_then(|v| v.data.as_mut()) {
            d.0 = 99;
        }
        assert_eq!(g.get(&"AA").and_then(|v| v.data.as_ref()), Some(&D(99, 1)));
    }

    #[test]
    fn test_vertex_edges() {
        let g = make_graph();
        assert_eq!(
            g.get(&"AA")
                .and_then(|v| Some(v.edges.iter().collect::<Vec<_>>())),
            Some(vec![&("BB", 1), &("CC", 2)])
        );
    }

    #[test]
    fn test_vertex_edges_mut() {
        let mut g = make_graph();
        g.get_mut(&"AA").and_then(|v| {
            v.edges.push(("ZZ", 99));
            Some(())
        });
        assert_eq!(
            g.get(&"AA")
                .and_then(|v| Some(v.edges.iter().collect::<Vec<_>>())),
            Some(vec![&("BB", 1), &("CC", 2), &("ZZ", 99)])
        );
    }
}
