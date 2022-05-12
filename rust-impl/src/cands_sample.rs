use std::collections::HashSet;
use std::fs;
use std::iter::zip;
use rand::seq::SliceRandom;

const WORD_LEN: usize = 5;
const STRAT_NUM: i32 = 3;

fn main() {
    // read in strategies
    let lines = match fs::read_to_string("../data/candidates.txt") {
        Ok(buff) => buff,
        Err(_) => panic!("Error reading candidates file."),
    };
    let lines: Vec<String> = lines.split("\n").map(String::from).collect();
    let strats: Vec<Vec<String>> = lines.iter().map(|x| x.split(",").map(String::from).collect()).collect();

    // read in words
    let db = Database::from_file("../data/words.txt");

    for i in 0..10{
        let guess_words = &strats[i];
        let sample: Vec<_> = db.words
            .choose_multiple(&mut rand::thread_rng(), 1_000)
            .collect();

        for hidden_word in sample {
            let mut info = Information::new();
            info.update_info(hidden_word, &guess_words);
            _ = db.query(info);
        }
        
    }
}

struct Information {
    has_set: HashSet<char>,
    not_set: HashSet<char>,
    known_pos: [char; WORD_LEN],
    not_pos: Vec<HashSet<char>>,
}

impl Information {
    pub fn new() -> Self {
        let mut empty_char_set: Vec<HashSet<char>> = vec![];
        for _ in 0..WORD_LEN {
            empty_char_set.push(HashSet::new())
        }
        return Self {
            has_set: HashSet::new(),
            not_set: HashSet::new(),
            known_pos: ['_'; WORD_LEN],
            not_pos: empty_char_set,
        };
    }

    pub fn update_info(&mut self, hidden_word: &String, guesses: &Vec<String>) {
        let hidden_set: HashSet<char> = hidden_word.chars().collect();
        for guess_word in guesses {
            for (i, (guess_char, hidden_char)) in
                zip(guess_word.chars(), hidden_word.chars()).enumerate()
            {
                // set tests
                if hidden_set.contains(&guess_char) {
                    self.has_set.insert(guess_char);
                } else {
                    self.not_set.insert(guess_char);
                }

                // position tests
                if guess_char == hidden_char {
                    self.known_pos[i] = guess_char;
                } else {
                    self.not_pos[i].insert(guess_char);
                }
            }
        }
    }

    /*
        TODO: see if there's way to make function that handles
        a generic type and calls the chain clone().into_iter().collect()
        to clean up this code below
    */
    pub fn print(&self) {
        println!("has set: ");
        let tmp: String = self.has_set.clone().into_iter().collect();
        println!("  {}", tmp);

        println!("not set:");
        let tmp: String = self.not_set.clone().into_iter().collect();
        println!("  {}", tmp);

        println!("known pos:");
        let tmp: String = self.known_pos.clone().into_iter().collect();
        println!("  {}", tmp);

        println!("not pos:");
        for v in self.not_pos.clone() {
            let tmp: String = v.into_iter().collect();
            println!("  {}", tmp);
        }
    }
}

struct Database {
    words: Vec<String>,
}

impl Database {
    pub fn new() -> Self {
        return Self { words: vec![] };
    }

    pub fn from_file(fp: &str) -> Self {
        let lines = match fs::read_to_string(fp) {
            Ok(buff) => buff,
            Err(_) => panic!("Error reading db file."),
        };
        return Self {
            words: lines.split("\n").map(String::from).collect(),
        };
    }

    pub fn query(&self, info: Information) -> Vec<String> {
        let mut possible: Vec<String> = vec![];

        for word in &self.words {
            let mut skip = false;
            for l in word.chars() {
                if info.not_set.contains(&l) {
                    skip = true;
                    break;
                }
            }
            if skip {
                continue;
            }

            let char_set: HashSet<char> = HashSet::from_iter(word.chars());
            if !info.has_set.is_subset(&char_set) {
                continue;
            }

            let mut no_skip = true;
            for (k, l) in zip(info.known_pos, word.chars()) {
                if k != '_' && k != l {
                    no_skip = false;
                    break;
                }
            }
            if !no_skip {
                continue;
            }

            skip = false;
            for (not_set, l) in zip(&info.not_pos, word.chars()) {
                if not_set.len() > 0 && not_set.contains(&l) {
                    skip = true;
                    break;
                }
            }
            if skip {
                continue;
            }

            possible.push(word.to_string());
        }

        return possible;
    }
}
