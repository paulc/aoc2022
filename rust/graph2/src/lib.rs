use std::cmp::Ord;
use std::collections::BTreeMap;
use std::fmt::Debug;
use std::fmt::Display;
use std::hash::Hash;

#[derive(Debug, Clone)]
pub struct E(());

impl Display for E {
    fn fmt(&self, _f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        Ok(())
    }
}

#[derive(Debug)]
pub struct Vertex<I, D>(I, Option<D>, Vec<(I, f64)>)
where
    I: Display + Clone + Eq + PartialEq + Ord + Hash,
    D: Display;

impl<I, D> Vertex<I, D>
where
    I: Display + Clone + Eq + PartialEq + Ord + Hash,
    D: Display,
{
    pub fn new(id: I, data: Option<D>, edges: Vec<(I, f64)>) -> Self {
        Self(id, data, edges)
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
    pub fn edges(&mut self) -> impl Iterator<Item = (I, f64)> + '_ {
        self.2.iter().cloned()
    }
    pub fn edges_mut(&mut self) -> &mut Vec<(I, f64)> {
        &mut self.2
    }
    pub fn add_edge(&mut self, to: I, cost: f64) {
        self.2.push((to, cost))
    }
}

#[derive(Debug)]
pub struct Graph<I, D>(BTreeMap<I, Vertex<I, D>>)
where
    I: Display + Clone + Eq + PartialEq + Ord + Hash,
    D: Display;

impl<I, D> Graph<I, D>
where
    I: Display + Clone + Eq + PartialEq + Ord + Hash,
    D: Display,
{
    pub fn new() -> Self {
        Self(BTreeMap::new())
    }
    pub fn new_from_edges(edges: Vec<(I, I, f64)>) -> Self {
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
}

impl<I, D> Display for Vertex<I, D>
where
    I: Display + Clone + Eq + PartialEq + Ord + Hash,
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
