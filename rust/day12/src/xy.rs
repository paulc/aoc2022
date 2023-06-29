#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash, PartialOrd, Ord)]
pub struct XY {
    pub x: i32,
    pub y: i32,
}

impl XY {
    pub fn new(x: i32, y: i32) -> XY {
        XY { x, y }
    }
    pub fn adjacent(&self) -> [XY; 4] {
        [
            XY::new(self.x, self.y - 1),
            XY::new(self.x + 1, self.y),
            XY::new(self.x, self.y + 1),
            XY::new(self.x - 1, self.y),
        ]
    }
}

impl std::fmt::Display for XY {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "({},{})", self.x, self.y)
    }
}
