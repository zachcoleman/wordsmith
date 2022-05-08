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
            skip = False
            for let in word:
                if let in info_state.not_set:
                    skip = True
                    break
            if skip:
                continue

            if not info_state.has_set.issubset(set(word)):
                continue

            no_skip = True
            for k, l in zip(info_state.known_pos, word):
                if k is not None and k != l:
                    no_skip = False
                    break
            if not no_skip:
                continue

            skip = False
            for not_set, l in zip(info_state.not_pos, word):
                if len(not_set) > 0 and l in not_set:
                    skip = True
                    break
            if skip:
                continue

            possible.append(word)

        return possible
