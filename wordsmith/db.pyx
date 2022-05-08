# cython: language_level=3

import os
from typing import List

from wordsmith.info import Information


cdef class WordDatabase:
    
    cdef public list words
    
    def __cinit__(self, fp: os.PathLike) -> None:
        with open(fp, "r") as f:
            self.words = f.read().splitlines()

    def query(self, info_state: Information) -> List[str]:        
        cdef list possible = []
        cdef str word
        cdef str l, k
        cdef bint skip, no_skip
        cdef set not_set

        for word in self.words:
            skip = False
            for l in word:
                if l in info_state.not_set:
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
