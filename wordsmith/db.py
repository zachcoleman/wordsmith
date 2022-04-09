import os
from typing import List

from wordsmith.info import Information


class WordDatabase:
    def __init__(self, fp: os.PathLike) -> None:
        with open(fp, "r") as f:
            self.words = f.read().splitlines()

    def query(self, info_state: Information) -> List[str]:
        possible = []
        for word in self.words:
            if not info_state.has_set.issubset(set(word)):
                continue

            bad_letters = False
            for letter in word:
                if letter in info_state.not_set:
                    bad_letters = True
                    break
            if bad_letters:
                continue

            match_letters = True
            for k, l in zip(info_state.known_pos, word):
                if k is not None and k != l:
                    match_letters = False
                    break
            if not match_letters:
                continue

            match_not_here = True
            for not_set, l in zip(info_state.not_pos, word):
                if len(not_set) > 0 and l in not_set:
                    match_not_here = False
                    break
            if not match_not_here:
                continue

            possible.append(word)

        return possible
