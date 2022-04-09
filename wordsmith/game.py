from typing import List

from .info import Information


class Game:
    def __init__(self, hidden_word: str) -> None:
        self.hidden_word = hidden_word

    def generate_information(self, guesses: List[str]) -> Information:
        ret_info = Information()
        for guess in guesses:
            for pos, letter in enumerate(guess):
                if letter in set(self.hidden_word):
                    if letter == self.hidden_word[pos]:
                        ret_info.known_pos[pos] = letter
                    else:
                        ret_info.not_pos[pos].add(letter)
                    ret_info.has_set.add(letter)
                else:
                    ret_info.not_set.add(letter)
        return ret_info
