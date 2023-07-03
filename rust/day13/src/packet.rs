use std::fmt::Display;

#[derive(Debug)]
pub struct PacketBuf<T>(Vec<Packet<T>>);

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
    fn fmt_packet(&self, root: PacketRef) -> String {
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
pub struct PacketRef(usize);

#[derive(Debug)]
pub enum Packet<T> {
    List(Vec<PacketRef>, Option<PacketRef>),
    Value(T),
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
        println!(
            "{}",
            PacketBuf::try_from("[1,2,3,[99,100],[[[10]]]]").unwrap()
        );
        println!("{}", PacketBuf::try_from("[[[]]]").unwrap());
    }
}
