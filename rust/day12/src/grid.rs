use crate::xy::XY;

#[derive(Debug)]
pub struct Grid<T>(Vec<Vec<T>>);

impl<T> Grid<T> {
    pub fn new(from: Vec<Vec<T>>) -> Self {
        Self(from)
    }
    pub fn check_bounds(&self, p: XY) -> bool {
        (p.y >= 0 && p.y < self.0.len() as i32)
            && ((self.0.len() > 0) && (p.x >= 0 && p.x < self.0[0].len() as i32))
    }
    pub fn get(&self, p: XY) -> Option<&T> {
        match self.check_bounds(p) {
            true => Some(&self.0[p.y as usize][p.x as usize]),
            false => None,
        }
    }
    pub fn adjacent(&self, p: XY) -> Vec<(&T, XY)> {
        p.adjacent()
            .into_iter()
            .filter(|&p| self.check_bounds(p))
            .map(|p| (&self.0[p.y as usize][p.x as usize], p))
            .collect()
    }
    pub fn iter(&self) -> impl Iterator<Item = (&T, XY)> {
        self.0.iter().enumerate().flat_map(|(y, r)| {
            r.iter()
                .enumerate()
                .map(move |(x, c)| (c, XY::new(x as i32, y as i32)))
        })
    }
}

impl<T: std::fmt::Display> std::fmt::Display for Grid<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for y in 0..self.0.len() {
            writeln!(
                f,
                "{}",
                self.0[y]
                    .iter()
                    .map(|e| format!("{e:>2}"))
                    .collect::<Vec<String>>()
                    .join(" ")
            )?
        }
        Ok(())
    }
}
