use std::fmt::Display;
use std::ops::{Add, Sub};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Point {
    pub x: i32,
    pub y: i32,
}

impl Point {
    pub fn new(x: i32, y: i32) -> Point {
        Point { x, y }
    }
    pub fn adjacent(&self) -> impl Iterator<Item = Point> + '_ {
        ADJACENT.into_iter().map(|offset| self.add(offset))
    }
}

impl Add<Offset> for Point {
    type Output = Self;
    fn add(self, Offset { dx, dy }: Offset) -> Self {
        Point {
            x: self.x + dx,
            y: self.y + dy,
        }
    }
}

impl Sub for Point {
    type Output = Offset;
    fn sub(self, Point { x, y }: Point) -> Offset {
        Offset {
            dx: self.x - x,
            dy: self.y - y,
        }
    }
}

impl TryFrom<&str> for Point {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        match s.split_once(',') {
            Some((x, y)) => match (x.parse::<i32>(), y.parse::<i32>()) {
                (Ok(x), Ok(y)) => Ok(Point { x, y }),
                _ => Err("Invalid input"),
            },
            _ => Err("Invalid input"),
        }
    }
}

impl Display for Point {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{},{}", self.x, self.y)
    }
}

#[derive(Debug, Clone, Copy, Eq, PartialEq)]
pub struct Offset {
    pub dx: i32,
    pub dy: i32,
}

pub const UP: Offset = Offset { dx: 0, dy: -1 };
pub const DOWN: Offset = Offset { dx: 0, dy: 1 };
pub const RIGHT: Offset = Offset { dx: 1, dy: 0 };
pub const LEFT: Offset = Offset { dx: -1, dy: 0 };

pub const ADJACENT: [Offset; 4] = [UP, RIGHT, DOWN, LEFT];

impl Offset {
    pub fn new(dx: i32, dy: i32) -> Offset {
        Offset { dx, dy }
    }
}

impl Add<Point> for Offset {
    type Output = Point;
    fn add(self, Point { x, y }: Point) -> Point {
        Point {
            x: x + self.dx,
            y: y + self.dy,
        }
    }
}

impl Add for Offset {
    type Output = Self;
    fn add(self, Offset { dx, dy }: Offset) -> Self {
        Offset {
            dx: dx + self.dx,
            dy: dy + self.dy,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_point() {
        let p1 = Point::new(1, 2);
        let o1 = Offset::new(3, 4);
        let p2 = p1 + o1;
        assert_eq!(p2, Point::new(4, 6));
        let o2 = p2 - p1;
        assert_eq!(o2, o1);
    }
    #[test]
    fn test_from() {
        assert_eq!(Point::try_from("-3,6"), Ok(Point::new(-3, 6)));
    }
    #[test]
    fn test_adjacent() {
        assert_eq!(
            Point::new(5, 5).adjacent().collect::<Vec<_>>(),
            vec![
                Point::new(5, 4),
                Point::new(6, 5),
                Point::new(5, 6),
                Point::new(4, 5)
            ]
        );
    }
}
