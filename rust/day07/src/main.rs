#![allow(unused)]

use std::cell::RefCell;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;
use std::rc::Rc;

type Out = usize;
type In = FS;

#[derive(Debug)]
struct ParseError;
impl std::error::Error for ParseError {}
impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Parse error")
    }
}

#[derive(Debug)]
struct Dir {
    id: usize,
    name: String,
    dirs: HashMap<String, usize>,
    files: HashMap<String, usize>,
    parent: Option<usize>,
    size: usize,
}

impl Dir {
    fn new(id: usize, name: String, parent: Option<usize>) -> Self {
        Dir {
            id,
            name,
            dirs: HashMap::new(),
            files: HashMap::new(),
            parent,
            size: 0,
        }
    }
}

#[derive(Debug)]
struct FS {
    dirs: Vec<Dir>,
    root: usize,
}

impl FS {
    fn new() -> Self {
        Self {
            dirs: vec![Dir::new(0, "/".to_string(), None)],
            root: 0,
        }
    }
    fn add_subdir(&mut self, cwd: usize, name: String) -> Result<(), String> {
        let id = self.dirs.len();
        let new = Dir::new(id, name.clone(), Some(cwd));
        self.dirs.push(new);
        let cwd = self.get_mut(cwd)?;
        cwd.dirs.insert(name, id);
        Ok(())
    }
    fn add_file(&mut self, cwd: usize, name: String, size: usize) -> Result<(), String> {
        let cwd = self.get_mut(cwd)?;
        cwd.files.insert(name, size);
        cwd.size += size;
        let mut parent = cwd.parent;
        while let Some(p) = parent {
            let p = self.get_mut(p)?; // self.dirs.get_mut(p).ok_or(format!("Cant get parent"))?;
            p.size += size;
            parent = p.parent;
        }
        Ok(())
    }
    fn chdir(&self, cwd: usize, name: &str) -> Result<usize, String> {
        match name {
            "/" => Ok(0),
            ".." => match self
                .dirs
                .get(cwd)
                .ok_or(format!("Cant get cwd {}", cwd))?
                .parent
            {
                Some(parent) => Ok(parent),
                None => Ok(cwd),
            },
            _ => match self
                .dirs
                .get(cwd)
                .ok_or(format!("Cant get cwd {}", cwd))?
                .dirs
                .get(name)
            {
                Some(next) => Ok(*next),
                None => Err(format!("No such dir: {}", name)),
            },
        }
    }
    fn get(&self, cwd: usize) -> Result<&Dir, String> {
        Ok(self.dirs.get(cwd).ok_or(format!("Cant get dir: {}", cwd))?)
    }
    fn get_mut(&mut self, cwd: usize) -> Result<&mut Dir, String> {
        Ok(self
            .dirs
            .get_mut(cwd)
            .ok_or(format!("Cant get dir: {}", cwd))?)
    }
    fn walk(&self, root: usize) -> Result<Vec<&Dir>, String> {
        let mut out: Vec<&Dir> = Vec::new();
        let cwd = self.get(root)?;
        out.push(self.get(root)?);
        for (k, v) in cwd.dirs.iter() {
            for d in self.walk(*v)? {
                out.push(d);
            }
        }
        Ok(out)
    }
}

fn parse_input(input: &mut impl Read) -> std::io::Result<In> {
    let mut fs = FS::new();
    let mut cwd: usize = fs.root;
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            match l.split_whitespace().collect::<Vec<&str>>().as_slice() {
                ["$", "cd", dir] => cwd = fs.chdir(cwd, dir).unwrap(), // pwd.chdir(&fs, dir).unwrap(),
                ["$", "ls"] => {}
                ["dir", dir] => fs.add_subdir(cwd, dir.to_string()).unwrap(),
                [size, name] => fs
                    .add_file(
                        cwd,
                        name.to_string(),
                        size.parse::<usize>()
                            .map_err(|e| Error::new(InvalidData, e))?,
                    )
                    .unwrap(),
                _ => return Err(Error::new(InvalidData, ParseError)),
            };
        }
    }
    Ok(fs)
}

fn part1(input: &In) -> Out {
    input
        .dirs
        .iter()
        .filter(|d| d.size < 100000)
        .map(|d| d.size)
        .sum::<usize>()
}

fn part2(input: &In) -> Out {
    let unused = 70000000 - input.get(input.root).unwrap().size;
    let required = 30000000 - unused;
    input
        .dirs
        .iter()
        .filter(|d| d.size > required)
        .map(|d| d.size)
        .collect::<Vec<usize>>()
        .iter()
        .min()
        .copied()
        .unwrap()
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let input = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&input));
    println!("Part2: {:?}", part2(&input));
    Ok(())
}

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&data), 95437);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&data), 24933642);
}

#[cfg(test)]
const TESTDATA: &str = "
$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k
";
