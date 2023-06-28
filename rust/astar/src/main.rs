pub mod graph;

use crate::graph::dijkstra::dijkstra;
use crate::graph::Graph;

fn main() {
    let mut g: Graph<String> = Graph::new();
    g.add_vertex("AAA".to_string());
    g.add_vertex("BBB".to_string());
    g.add_vertex("CCC".to_string());
    g.add_edge("AAA".to_string(), "BBB".to_string(), 1.0);
    g.add_edge("AAA".to_string(), "CCC".to_string(), 1.0);
    g.add_edge("BBB".to_string(), "DDD".to_string(), 2.0);
    g.add_edge("CCC".to_string(), "DDD".to_string(), 2.0);
    g.add_edge_symmetric("DDD".to_string(), "AAA".to_string(), 3.0);
    g.add_edge_symmetric("DDD".to_string(), "EEE".to_string(), 4.0);
    println!("{}", g);
}
