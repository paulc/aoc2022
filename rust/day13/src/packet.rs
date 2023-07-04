use std::fmt::Display;

#[derive(Debug)]
pub struct PacketBuf<T>(pub Vec<Packet<T>>);

impl<T> PacketBuf<T> {
    pub fn new() -> Self {
        Self(Vec::new())
    }
    pub fn get_ref(&mut self, p: Packet<T>) -> PacketRef {
        let n = self.0.len();
        self.0.push(p);
        PacketRef(n)
    }
    pub fn push(&mut self, parent: PacketRef, child: PacketRef) -> Result<(), &'static str> {
        let parent = self.0.get_mut(parent.0).ok_or_else(|| "Cant find parent")?;
        match parent {
            Packet::List(l, _) => {
                l.push(child);
                Ok(())
            }
            Packet::Value(_) => Err("Cant push to Packet::Value"),
        }
    }
    pub fn iter(&self, root: PacketRef) -> impl Iterator<Item = (PacketRef, &Packet<T>)> {
        match self.0.get(root.0) {
            Some(Packet::List(l, _)) => l.iter().map(|p| (*p, self.0.get(p.0).unwrap())),
            _ => panic!("Invalid root ref"),
        }
    }
    pub fn _iter(&self, root: PacketRef) -> PacketIter<T> {
        PacketIter {
            buf: self,
            current: root,
            i: 0,
        }
    }
}

#[derive(Debug)]
pub struct PacketIter<'a, T> {
    buf: &'a PacketBuf<T>,
    current: PacketRef,
    i: usize,
}

impl<'a, T> Iterator for PacketIter<'a, T> {
    type Item = (PacketRef, &'a Packet<T>);
    fn next(&mut self) -> Option<Self::Item> {
        if let Some(Packet::List(l, _)) = self.buf.0.get(self.current.0) {
            if let Some(next) = l.get(self.i) {
                self.i += 1;
                Some((*next, self.buf.0.get(next.0).unwrap()))
            } else {
                None
            }
        } else {
            None
        }
    }
}

impl TryFrom<&str> for PacketBuf<i32> {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        let mut buf: PacketBuf<i32> = PacketBuf::new();
        let b = s.as_bytes();
        let mut i = 0;
        let mut current: Option<PacketRef> = None;
        while i < b.len() {
            match b[i] {
                b'[' => {
                    let new = buf.get_ref(Packet::new_list(current));
                    if let Some(current) = current {
                        buf.push(current, new);
                    }
                    current = Some(new);
                    i += 1;
                }
                b']' => {
                    let p = buf
                        .0
                        .get_mut(current.unwrap().0) // XXX
                        .ok_or_else(|| "Cant find parent")?;
                    match p {
                        Packet::List(_, parent) => current = *parent,
                        Packet::Value(_) => return Err("Error closing list - invalid type"),
                    }
                    i += 1;
                }
                b'0'..=b'9' => {
                    let mut n = 0_i32;
                    while b[i].is_ascii_digit() {
                        n = (n * 10) + ((b[i] - b'0') as i32);
                        i += 1;
                    }
                    let new = buf.get_ref(Packet::new_value(n));
                    buf.push(current.unwrap(), new); // XXX
                }
                b',' => {
                    i += 1;
                }
                _ => return Err("Invalid input"),
            }
        }
        Ok(buf)
    }
}

impl<T: Display> PacketBuf<T> {
    pub fn fmt_packet(&self, root: PacketRef) -> String {
        if let Some(root) = self.0.get(root.0) {
            match root {
                Packet::List(l, _) => {
                    let mut out: Vec<String> = Vec::new();
                    for p in l {
                        out.push(self.fmt_packet(*p));
                    }
                    format!("[{}]", out.join(","))
                }
                Packet::Value(v) => v.to_string(),
            }
        } else {
            format!("Cant find Packet: {}", root.0)
        }
    }
}

impl<T: Display> Display for PacketBuf<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.fmt_packet(PacketRef(0))) // Assume root is packet 0
    }
}

#[derive(Debug, Clone, Copy)]
pub struct PacketRef(pub usize);

#[derive(Debug)]
pub enum Packet<T> {
    List(Vec<PacketRef>, Option<PacketRef>),
    Value(T),
}

impl<T: Display> Display for Packet<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Packet::List(l, _) => {
                write!(f, "[..{}]", l.len())
            }
            Packet::Value(v) => write!(f, "{}", v),
        }
    }
}

impl<T> Packet<T> {
    pub fn new_value(v: T) -> Self {
        Packet::Value(v)
    }
    pub fn new_list(parent: Option<PacketRef>) -> Self {
        Packet::List(Vec::new(), parent)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_packet() {
        let mut s: PacketBuf<i32> = PacketBuf::new();
        let mut root = s.get_ref(Packet::new_list(None));
        let mut current = root;
        for i in 0..3 {
            let p = s.get_ref(Packet::new_value(i));
            s.push(current, p);
        }
        for i in 0..3 {
            let p = s.get_ref(Packet::new_list(Some(current)));
            s.push(current, p);
            current = p;
            for i in 0..3 {
                let p = s.get_ref(Packet::new_value(i));
                s.push(current, p);
            }
        }
        for i in 0..3 {
            let p = s.get_ref(Packet::new_value(i));
            s.push(root, p);
        }
        assert_eq!(s.to_string(), "[0,1,2,[0,1,2,[0,1,2,[0,1,2]]],0,1,2]");
    }
    #[test]
    fn test_from() {
        for p in ["[1,2,3,[99,100],[[[10]]]]", "[1,1,3,1,1]", "[[1],[2,3,4]]"] {
            assert_eq!(PacketBuf::try_from(p).unwrap().to_string(), p);
        }
    }

    #[test]
    fn test_iter() {
        let buf = PacketBuf::try_from("[1,2,3,[99,100],[[[10]]]]").unwrap();
        assert_eq!(
            buf.iter(PacketRef(0)).map(|(i, _)| i.0).collect::<Vec<_>>(),
            vec![1, 2, 3, 4, 7]
        );
        assert_eq!(
            buf.iter(PacketRef(4)).map(|(i, _)| i.0).collect::<Vec<_>>(),
            vec![5, 6]
        );
        assert_eq!(
            buf.iter(PacketRef(7)).map(|(i, _)| i.0).collect::<Vec<_>>(),
            vec![8]
        );
    }
}
