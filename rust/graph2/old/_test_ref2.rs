use std::cell::RefCell;
use std::fmt::Debug;
use std::fmt::Display;
use std::hash::{Hash, Hasher};
use std::rc::Rc;

#[derive(Debug)]
pub struct VertexRef<T>(Rc<RefCell<Vertex<T>>>)
where
    T: Display + PartialEq + Eq + Hash;

impl<T> Eq for VertexRef<T> where T: Display + PartialEq + Eq + Hash {}

impl<T> PartialEq for VertexRef<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    fn eq(&self, other: &Self) -> bool {
        self.0.borrow().v == other.0.borrow().v
    }
}

impl<T> Hash for VertexRef<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.0.borrow().v.hash(state);
    }
}

impl<T> VertexRef<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    pub fn new(v: Vertex<T>) -> Self {
        Self(Rc::new(RefCell::new(v)))
    }
    pub fn clone(&self) -> Self {
        Self(Rc::clone(&self.0))
    }

    pub fn edges(&self) -> Vec<(VertexRef<T>, f64)> {
        self.0
            .borrow()
            .e
            .iter()
            .map(|(r, c)| (r.clone(), c.clone()))
            .collect()
    }
    pub fn add_edge(&mut self, v: &VertexRef<T>, c: f64) {
        self.0.borrow_mut().e.push((v.clone(), c))
    }
}

impl<T> Display for VertexRef<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.0.borrow())
    }
}

pub struct Vertex<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    v: T,
    e: Vec<(VertexRef<T>, f64)>,
}

impl<T> Vertex<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    pub fn new(v: T, e: Vec<(&VertexRef<T>, f64)>) -> Self {
        Self {
            v,
            e: e.iter().map(|&(v, c)| (v.clone(), c.clone())).collect(),
        }
    }
}

impl<T> Display for Vertex<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{} -> ", self.v)?;
        let mut n = self.e.len();
        for (d, c) in &self.e {
            let v = d.0.borrow();
            write!(f, "[{}]({})", v.v, c)?;
            n -= 1;
            if n > 0 {
                write!(f, ",")?;
            }
        }
        Ok(())
    }
}

impl<T> Debug for Vertex<T>
where
    T: Display + PartialEq + Eq + Hash,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Vertex {{ {} -> ", self.v)?;
        let mut n = self.e.len();
        for (d, c) in &self.e {
            let v = d.0.borrow();
            write!(f, "[{}]({})", v.v, c)?;
            n -= 1;
            if n > 0 {
                write!(f, ",")?;
            }
        }
        write!(f, " }}")?;
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    fn make_graph() -> HashMap<String, RefCell<VertexRef<String>>> {
        let mut g: HashMap<String, RefCell<VertexRef<String>>> = HashMap::new();
        for v in vec!["AA", "BB", "CC", "DD", "EE"] {
            g.insert(
                v.to_string(),
                RefCell::new(VertexRef::new(Vertex::new(v.to_string(), vec![]))),
            );
        }
        g["AA"].borrow_mut().add_edge(&g["BB"].borrow(), 1.0);
        g["AA"].borrow_mut().add_edge(&g["CC"].borrow(), 1.0);
        g
    }

    #[test]
    fn test_graph() {
        let g = make_graph();
        println!("{:#?}", g);
    }
}
