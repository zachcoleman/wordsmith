from typing import List, Set

from .constants import WORD_LEN


class Information:
    def __init__(
        self,
        known_pos: List[str] = None,
        not_pos: List[Set[str]] = None,
        has_set: Set[str] = None,
        not_set: Set[str] = None,
    ) -> None:
        self.known_pos = known_pos if known_pos else [None for _ in range(WORD_LEN)]
        self.not_pos = not_pos if not_pos else [set() for _ in range(WORD_LEN)]
        self.has_set = has_set if has_set else set()
        self.not_set = not_set if not_set else set()

    def __repr__(self) -> str:
        return "\n".join(
            [
                " ".join(["_" if let is None else let for let in self.known_pos]),
                "Letters not in position:",
                "\n".join(
                    [
                        " " + (str(list(not_set)) if k is None else k)
                        for k, not_set in zip(self.known_pos, self.not_pos)
                    ]
                ),
                "Has letters: " + ",".join([let for let in self.has_set]),
                "Does not have letter: " + ",".join([let for let in self.not_set]),
            ]
        )
