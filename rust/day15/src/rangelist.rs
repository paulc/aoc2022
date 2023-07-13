use std::cmp::Ordering;
use std::collections::BTreeSet;

#[derive(Debug)]
pub struct Rangelist<T: Clone + Ord>(BTreeSet<[T; 2]>);

impl<T: Clone + Ord> Rangelist<T> {
    pub fn new() -> Rangelist<T> {
        Rangelist(BTreeSet::new())
    }
    pub fn iter(&self) -> impl Iterator<Item = &[T; 2]> {
        self.0.iter()
    }
    pub fn add(&mut self, v: [T; 2]) {
        assert!(v[0] <= v[1]);
        self.0.insert(v);
    }
    pub fn first(&self) -> Option<&[T; 2]> {
        self.0.first()
    }
    pub fn pop_first(&mut self) -> Option<[T; 2]> {
        self.0.pop_first()
    }
    pub fn len(&self) -> usize {
        self.0.len()
    }
    pub fn contains(&self, v: T) -> bool {
        for r in &self.0 {
            if v >= r[0] && v <= r[1] {
                return true;
            }
        }
        false
    }
    pub fn coalesce(&mut self) {
        let mut i = self.0.iter();
        let mut out: BTreeSet<[T; 2]> = BTreeSet::new();
        match i.next() {
            None => {}
            Some(start) => {
                let mut current = start.clone();
                while let Some(next) = i.next() {
                    match (
                        current[0].cmp(&next[0]),
                        current[1].cmp(&next[0]),
                        current[1].cmp(&next[1]),
                    ) {
                        // c0---c1 n0---n1
                        (Ordering::Less, Ordering::Less, _) => {
                            out.insert(current);
                            current = next.clone();
                        }
                        // c0---c1
                        //    n0---n1
                        (
                            Ordering::Less | Ordering::Equal,
                            Ordering::Greater | Ordering::Equal,
                            Ordering::Less | Ordering::Equal,
                        ) => {
                            current = [current[0].clone(), next[1].clone()];
                        }
                        // c0--------c1
                        //    n0---n1
                        (Ordering::Less | Ordering::Equal, _, Ordering::Greater) => {
                            current = [current[0].clone(), current[1].clone()];
                        }
                        _ => panic!("Shouldnt get here"),
                    }
                }
                out.insert(current);
                self.0 = out;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn new_rangelist() -> Rangelist<i32> {
        let mut rl: Rangelist<i32> = Rangelist::new();
        rl.add([7, 10]);
        rl.add([2, 5]);
        rl.add([1, 2]);
        rl.add([-5, -3]);
        rl.add([-10, -2]);
        return rl;
    }

    #[test]
    fn test_first() {
        let rl = new_rangelist();
        assert_eq!(rl.first(), Some(&[-10, -2]));
    }

    #[test]
    fn test_contains() {
        let rl = new_rangelist();
        assert_eq!(rl.len(), 5);
        assert_eq!(rl.contains(3), true);
        assert_eq!(rl.contains(6), false);
        assert_eq!(rl.contains(-7), true);
        assert_eq!(rl.contains(-11), false);
    }

    #[test]
    #[should_panic]
    fn test_invalid() {
        let mut rl = new_rangelist();
        rl.add([5, 0]);
    }

    #[test]
    fn test_iter() {
        let rl = new_rangelist();
        assert_eq!(
            rl.iter().cloned().collect::<Vec<_>>(),
            vec![[-10, -2], [-5, -3], [1, 2], [2, 5], [7, 10]]
        );
    }

    #[test]
    fn test_coalsece() {
        let mut rl = new_rangelist();
        rl.coalesce();
        assert_eq!(
            rl.iter().cloned().collect::<Vec<_>>(),
            vec![[-10, -2], [1, 5], [7, 10]]
        );
    }
}
