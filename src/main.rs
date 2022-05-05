use std::collections::HashSet;
use std::fs;
use std::iter::zip;

const WORD_LEN: usize = 5;
const STRAT_NUM: i32 = 3;
const EMPTY_CHAR_VEC: Vec<char> = vec![];


fn main() {
    let lines = match fs::read_to_string("./data/candidates.txt") {
        Ok(buff) => buff,
        Err(_) => panic!("Error reading candidates file."),
    };

    let hidden_word = "hairy";
    let guess_words = vec!["adieu", "chair", "wilks"];

    let mut info = Information::new();
    info.update_info(hidden_word, guess_words);
    info.print();
}

struct Information {
    has_set: HashSet<char>,
    not_set: HashSet<char>,
    known_pos: [char; WORD_LEN],
    not_pos: [Vec<char>; WORD_LEN],
}

impl Information {
    pub fn new() -> Self {
        return Self {
            has_set: HashSet::new(),
            not_set: HashSet::new(),
            known_pos: ['_'; WORD_LEN],
            not_pos: [EMPTY_CHAR_VEC; WORD_LEN],
        };
    }

    pub fn update_info(&mut self, hidden_word: &str, guesses: Vec<&str>) {
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
                    self.not_pos[i].push(guess_char);
                }
            }
        }
    }

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
