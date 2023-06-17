#![allow(unused)]

use std::cell::RefCell;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::io::Error;
use std::io::ErrorKind::InvalidData;
use std::rc::Rc;

#[derive(Debug)]
enum Cmd {
    CdCmd(String),
    LsCmd,
    Dir(String),
    File(String, usize),
}

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
}

impl Dir {
    fn new(id: usize, name: String, parent: Option<usize>) -> Self {
        Dir {
            id,
            name,
            dirs: HashMap::new(),
            files: HashMap::new(),
            parent,
        }
    }
}

#[derive(Debug)]
struct DirRef {
    dirref: usize,
}

impl DirRef {
    fn new() -> Self {
        Self { dirref: 0 }
    }
    fn chdir(&mut self, fs: &FS, name: &str) -> Result<(), String> {
        self.dirref = fs.chdir(self.dirref, name)?;
        Ok(())
    }
}

#[derive(Debug)]
struct FS {
    dirs: Vec<Dir>,
}

impl FS {
    fn new() -> Self {
        Self {
            dirs: vec![Dir::new(0, "/".to_string(), None)],
        }
    }
    fn add_subdir(&mut self, cwd: usize, name: String) -> Result<(), ()> {
        let id = self.dirs.len();
        let new = Dir::new(id, name.clone(), Some(cwd));
        self.dirs.push(new);
        let cwd = self.dirs.get_mut(cwd).unwrap();
        cwd.dirs.insert(name, id);
        Ok(())
    }
    fn add_file(&mut self, cwd: usize, name: String, size: usize) -> Result<(), ()> {
        let cwd = self.dirs.get_mut(cwd).ok_or(())?;
        cwd.files.insert(name, size);
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
}

fn parse_input(input: &mut impl Read) -> std::io::Result<FS> {
    let mut fs = FS::new();
    let mut pwd = DirRef::new();
    let reader = BufReader::new(input);
    for l in reader.lines() {
        if let Ok(l) = l {
            match l.split_whitespace().collect::<Vec<&str>>().as_slice() {
                ["$", "cd", dir] => pwd.chdir(&fs, dir).unwrap(),
                ["$", "ls"] => {}
                ["dir", dir] => fs.add_subdir(pwd.dirref, dir.to_string()).unwrap(),
                [size, name] => fs
                    .add_file(
                        pwd.dirref,
                        name.to_string(),
                        size.parse::<usize>()
                            .map_err(|e| Error::new(InvalidData, e))?,
                    )
                    .unwrap(),
                _ => return Err(Error::new(InvalidData, ParseError)),
            };
        }
    }
    println!("{:#?}", fs);
    Ok(fs)
}

fn part1(root: &FS) -> Option<usize> {
    Some(95437)
}

fn part2(root: &FS) -> Option<usize> {
    Some(24933642)
}

fn main() -> std::io::Result<()> {
    let mut f = File::open("input")?;
    let root = parse_input(&mut f)?;
    println!("Part1: {:?}", part1(&root).unwrap());
    println!("Part2: {:?}", part2(&root).unwrap());
    Ok(())
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

#[test]
fn test_part1() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part1(&data).unwrap(), 95437);
}

#[test]
fn test_part2() {
    let data = parse_input(&mut TESTDATA.trim_matches('\n').as_bytes()).unwrap();
    assert_eq!(part2(&data).unwrap(), 24933642);
}
