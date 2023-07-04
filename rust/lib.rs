use std::fmt::Display;

#[derive(Debug)]
pub struct XX<T> {}

#[derive(Debug)]
pub struct XX_Iter<'a, T> {
    buf: &'a XX<T>,
}

impl<'a, T> Iterator for XX_Iter<'a, T> {
    type Item = ();
    fn next(&mut self) -> Option<Self::Item> {
        None
    }
}

impl TryFrom<&str> for XX<i32> {
    type Error = &'static str;
    fn try_from(s: &str) -> Result<Self, Self::Error> {
        Err("Invalid input")
    }
}

impl<T: Display> Display for XX<T> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", "")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_() {}
}
