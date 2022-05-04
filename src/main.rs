use std::collections::HashSet;
use std::fs;
use std::iter::zip;

const WORD_LEN: usize = 5;
const STRAT_NUM: i32 = 3;

fn main() {
    let lines = match fs::read_to_string("./data/candidates.txt") {
        Ok(buff) => buff,
        Err(_) => panic!("Error reading candidates file."),
    };

    let hidden_word = "hairy";
    let guess_words = vec!["adieu", "chair", "wilks"];

    const EMPTY_CHAR_VEC: Vec<char> = vec![];
    let hidden_set: HashSet<char> = hidden_word.chars().collect();

    let mut has_set: HashSet<char> = HashSet::new();
    let mut not_set: HashSet<char> = HashSet::new();
    let mut known_pos = ['_'; WORD_LEN];
    let mut not_pos: [Vec<char>; WORD_LEN] = [EMPTY_CHAR_VEC; WORD_LEN];

    for guess_word in guess_words {
        for (i, (guess_char, hidden_char)) in
            zip(guess_word.chars(), hidden_word.chars()).enumerate()
        {
            // set tests
            if hidden_set.contains(&guess_char) {
                has_set.insert(guess_char);
            } else {
                not_set.insert(guess_char);
            }

            // position tests
            if guess_char == hidden_char {
                known_pos[i] = guess_char;
            } else {
                not_pos[i].push(guess_char);
            }
        }
    }

    println!("has set: ");
    {
        let tmp: String = has_set.into_iter().collect();
        println!("  {}", tmp);
    }

    println!("not set:");
    {
        let tmp: String = not_set.into_iter().collect();
        println!("  {}", tmp);
    }

    println!("known pos:");
    {
        let tmp: String = known_pos.into_iter().collect();
        println!("  {}", tmp);
    }

    println!("not pos:");
    for v in not_pos {
        let tmp: String = v.into_iter().collect();
        println!("  {}", tmp);
    }
}
