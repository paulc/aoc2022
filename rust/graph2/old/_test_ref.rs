#[derive(Debug, Clone)]
pub struct Edge<T>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
{
    to: T,
    cost: f64,
}

impl<T> Edge<T>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
{
    pub fn new(to: T, cost: f64) -> Self {
        Self { to, cost }
    }
}

impl<I> Display for Edge<I>
where
    I: Clone + Copy + Eq + PartialEq + Ord + Hash + Display,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "[{}]({})", self.to, self.cost)
    }
}

#[derive(Debug)]
pub struct Vertex<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
    D: Clone,
{
    node: T,
    edges: Vec<Edge<T>>,
    data: Option<D>,
    r: GraphRef<T, D>,
}

impl<T, D> Vertex<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
    D: Clone,
{
    pub fn new(node: T, edges: Vec<Edge<T>>, data: Option<D>, r: GraphRef<T, D>) -> Self {
        Self {
            node,
            edges,
            data,
            r,
        }
    }
    pub fn add_edge(&mut self, e: Edge<T>) {
        self.edges.push(e)
    }
    pub fn clone(&self) -> Self {
        Self {
            node: self.node.clone(),
            edges: self.edges.clone(),
            data: self.data.clone(),
            r: self.r.clone(),
        }
    }
}

#[derive(Debug)]
pub struct Graph<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
    D: Clone,
{
    g: BTreeMap<T, Vertex<T, D>>,
}

#[derive(Debug)]
pub struct GraphRef<T, D>(Rc<RefCell<Graph<T, D>>>)
where
    T: Clone + Copy + Eq + PartialEq + Ord + Ord + Hash,
    D: Clone;

impl<T, D> GraphRef<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Ord + Hash,
    D: Clone,
{
    pub fn new() -> Self {
        Self(Rc::new(RefCell::new(Graph::new())))
    }

    pub fn new_from_edges(edges: Vec<(T, T, f64)>) -> Self {
        let out = Self::new();
        for (v1, v2, cost) in edges {
            out.0
                .borrow_mut()
                .g
                .entry(v2)
                .or_insert_with(|| Vertex::new(v2, vec![], None, out.clone()));
            out.0
                .borrow_mut()
                .g
                .entry(v1)
                .or_insert_with(|| Vertex::new(v1, vec![], None, out.clone()))
                .add_edge(Edge::new(v2, cost));
        }
        out
    }

    pub fn clone(&self) -> Self {
        Self(Rc::clone(&self.0))
    }

    pub fn get(&self, id: &T) -> Option<Vertex<T, D>> {
        self.0.borrow().g.get(id).map(|v| v.clone())
    }

    pub fn apply<F, G>(&self, id: &T, f: F) -> Option<G>
    where
        F: Fn(&Vertex<T, D>) -> G,
    {
        self.0.borrow().g.get(id).and_then(|v| Some(f(v)))
    }

    pub fn update<F, G>(&mut self, id: &T, mut f: F) -> Option<G>
    where
        F: FnMut(&mut Vertex<T, D>) -> G,
    {
        self.0.borrow_mut().g.get_mut(id).and_then(|v| Some(f(v)))
    }

    pub fn map<F, G>(&self, f: F) -> Vec<G>
    where
        F: Fn(&Vertex<T, D>) -> G,
    {
        self.0
            .borrow()
            .g
            .iter()
            .map(|(_, v)| f(v))
            .collect::<Vec<G>>()
    }

    pub fn add_vertex(&mut self, v: T, e: Vec<Edge<T>>, d: Option<D>) {
        let mut binding = self.0.borrow_mut();
        let mut v = binding
            .g
            .entry(v)
            .or_insert_with(|| Vertex::new(v, vec![], None, self.clone()));
        e.into_iter().for_each(|e| v.add_edge(e));
        if let Some(_) = d {
            v.data = d
        }
    }
}

impl<T, D> Graph<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + Hash,
    D: Clone,
{
    pub fn new() -> Self {
        Graph { g: BTreeMap::new() }
    }
    /*
    fn add_symmetric_edges(&mut self) -> Result<(), &'static str> {
        let mut symmetric: HashMap<I, Vec<Edge<I>>> = HashMap::new();
        for (i, v) in &self.g {
            for e in &v.edges {
                symmetric
                    .entry(e.to)
                    .or_insert(Vec::new())
                    .push(Edge::new(*i, e.cost));
            }
        }
        for (i, v) in symmetric.iter_mut() {
            let mut rev = self.get_mut(&i).ok_or_else(|| "Reverse vertex not found")?;
            rev.edges.append(v);
        }
        Ok(())
    }
    */
}

fn fmt_edges<T, D>(v: &Vertex<T, D>) -> String
where
    T: Clone + Copy + Eq + PartialEq + Ord + PartialOrd + Hash + Display,
    D: Clone,
{
    format!(
        "{}",
        v.edges
            .iter()
            .map(|e| e.to_string())
            .collect::<Vec<_>>()
            .join(",")
    )
}

struct DisplayGraph<T, D>(GraphRef<T, D>)
where
    T: Clone + Copy + Eq + PartialEq + Ord + PartialOrd + Hash + Display,
    D: Display + Clone;

impl<T, D> Display for DisplayGraph<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + PartialOrd + Hash + Display,
    D: Display + Clone,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for (_, v) in self.0 .0.borrow().g.iter() {
            writeln!(
                f,
                "{} -> {} [{}]",
                v.node,
                fmt_edges(v),
                v.data.as_ref().map_or("".to_string(), |x| x.to_string()),
            )?;
        }
        Ok(())
    }
}

impl<T, D> Display for GraphRef<T, D>
where
    T: Clone + Copy + Eq + PartialEq + Ord + PartialOrd + Hash + Display + std::fmt::Debug,
    D: Clone + std::fmt::Debug,
{
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for (_, v) in self.0.borrow().g.iter() {
            writeln!(f, "{} -> {}", v.node, fmt_edges(v))?;
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

    impl Display for V {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "[{}]", self.id)
        }
    }

    fn make_graph() -> GraphRef<&'static str, i64> {
        GraphRef::new_from_edges(vec![
            ("AA", "BB", 1.0),
            ("AA", "CC", 2.0),
            ("BB", "DD", 3.0),
            ("CC", "DD", 1.0),
            ("DD", "EE", 1.0),
        ])
    }

    #[test]
    fn test_graph() {
        let mut g: GraphRef<i64, ()> = GraphRef::new();
        g.add_vertex(
            1,
            vec![Edge::new(2, 1.0), Edge::new(3, 1.0), Edge::new(4, 1.0)],
            None,
        );
        g.add_vertex(2, vec![], None);
        g.add_vertex(3, vec![Edge::new(4, 1.0)], None);
        println!("{}", g);
    }

    #[test]
    fn test_graph_from_edges() {
        let g = make_graph();
        println!("{}", g);
    }

    #[test]
    fn test_graph_update_vertex() {
        let mut g = make_graph();
        assert_eq!(
            g.update(&"AA", |v| {
                v.data = Some(9999);
                ()
            }),
            Some(())
        );
        assert_eq!(
            g.update(&"ZZ", |v| {
                v.data = Some(9999);
            }),
            None
        );
        assert_eq!(g.get(&"AA").unwrap().data, Some(9999));
        println!("{}", g);
    }

    #[test]
    fn test_graph_map() {
        let g = make_graph();
        let m = g.map(|v| (v.node, v.data));
        println!("{:?}", m);
    }

    #[test]
    fn test_graph_ref() {
        let g = make_graph();
        let v = g.get(&"AA").unwrap();
        let mut g2 = v.r;
        g2.update(&"EE", |v| {
            v.data = Some(9999);
        })
        .unwrap();
        println!("{}", DisplayGraph(g));
    }

    #[test]
    fn test_displaygraph() {
        let mut g: GraphRef<i64, &str> = GraphRef::new();
        g.add_vertex(
            1,
            vec![Edge::new(2, 1.0), Edge::new(3, 1.0), Edge::new(4, 1.0)],
            Some("AA"),
        );
        g.add_vertex(2, vec![], Some("BB"));
        g.add_vertex(3, vec![], Some("CC"));
        println!("{}", DisplayGraph(g));
    }
}
