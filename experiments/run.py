from concurrent.futures import ProcessPoolExecutor
from functools import partial
from typing import List

import numpy as np
import pandas as pd

from wordsmith.db import WordDatabase
from wordsmith.game import Game


def eval_strat(db: WordDatabase, strat: List[str]):
    rs = np.random.choice(db.words, 1_000)
    words_left = []
    for word in rs:
        game = Game(word)
        info = game.generate_information(strat)
        words_left.append(len(db.query(info)))

    return {
        **{f"w{i}": w for i, w in enumerate(strat)},
        "median": np.median(words_left),
        "mean": np.mean(words_left),
    }


if __name__ == "__main__":
    db = WordDatabase("../data/words.txt")
    with open("../data/candidates.txt") as f:
        strats = []
        for line in f.readlines():
            strats.append(line[:-1].split(","))

    indices = np.random.choice(range(len(strats)), 1_000)
    strats_sample = [strats[ix] for ix in indices]

    with ProcessPoolExecutor() as ex:
        data = ex.map(partial(eval_strat, db), strats_sample)

    pd.DataFrame(data).to_csv("results.csv")
