use std::cmp::Ord;
use std::cmp::PartialOrd;
use std::collections::HashMap;
use std::fmt::Display;
use std::hash::Hash;

trait Id<I>
where
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    fn id(&self) -> I;
}

#[derive(Debug)]
struct Edge<I>
where
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    to: I,
    cost: f64,
}

impl<I> Edge<I>
where
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    fn new(to: I, cost: f64) -> Self {
        Self { to, cost }
    }
}

impl<I> Display for Edge<I>
where
    I: Clone + Copy + Eq + PartialEq + Hash + Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "id={}({})", self.to, self.cost)
    }
}

#[derive(Debug)]
struct Vertex<T, I>
where
    T: Id<I>,
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    node: T,
    edges: Vec<Edge<I>>,
}

impl<T, I> Vertex<T, I>
where
    T: Id<I>,
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    fn new(node: T, edges: Vec<Edge<I>>) -> Self {
        Self { node, edges }
    }
    fn add_edge(&mut self, e: Edge<I>) {
        self.edges.push(e)
    }
}

#[derive(Debug)]
struct Graph<T, I>
where
    T: Id<I>,
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    v: HashMap<I, Vertex<T, I>>,
}

impl<T, I> Graph<T, I>
where
    T: Id<I>,
    I: Clone + Copy + Eq + PartialEq + Hash,
{
    fn new() -> Self {
        Graph { v: HashMap::new() }
    }
    fn get(&self, id: &I) -> Option<&Vertex<T, I>> {
        self.v.get(id)
    }
    fn get_mut(&mut self, id: &I) -> Option<&mut Vertex<T, I>> {
        self.v.get_mut(id)
    }
    fn iter(&self) -> impl Iterator<Item = &Vertex<T, I>> {
        self.v.iter().map(|(k, v)| v)
    }
    fn add_vertex(&mut self, v: T, e: Vec<Edge<I>>) {
        self.v.insert(v.id(), Vertex::new(v, e));
    }
}

impl<T, I> Display for Graph<T, I>
where
    T: Id<I> + Display,
    I: Clone + Copy + Eq + PartialEq + Ord + PartialOrd + Hash + Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let mut keys = self.v.iter().map(|(k, v)| k).collect::<Vec<_>>();
        keys.sort();
        for k in keys {
            let v = self.get(k).expect("vertex");
            writeln!(
                f,
                "{} {}",
                v.node,
                v.edges
                    .iter()
                    .map(|e| e.to_string())
                    .collect::<Vec<_>>()
                    .join(",")
            )?;
        }
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[derive(Debug)]
    struct V {
        id: i64,
    }

    impl Id<i64> for V {
        fn id(&self) -> i64 {
            self.id
        }
    }

    impl Display for V {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "[id={}]", self.id)
        }
    }

    #[test]
    fn test_graph() {
        let mut g: Graph<V, i64> = Graph::new();
        g.add_vertex(V { id: 1 }, vec![Edge::new(2, 1.0), Edge::new(3, 1.0)]);
        g.add_vertex(V { id: 2 }, vec![]);
        g.add_vertex(V { id: 3 }, vec![]);
        println!("{}", g);
    }
}
