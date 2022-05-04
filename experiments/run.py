import numpy as np
import pandas as pd

from wordsmith.db import WordDatabase
from wordsmith.game import Game

# times:
# - original code: ~34s
# - python + optimal reordering: ~15s
# - typing cython + reordering optimal: ~7.3s

if __name__ == "__main__":
    db = WordDatabase("../data/words.txt")
    with open("../data/candidates.txt") as f:
        strats = []
        for line in f.readlines():
            strats.append(line[:-1].split(","))

    data = []
    for ix in range(10):
        strat = strats[ix]
        rs = np.random.choice(db.words, 1_000)
        words_left = []
        for word in rs:
            game = Game(word)
            info = game.generate_information(strat)
            words_left.append(len(db.query(info)))

        data.append(
            {
                **{f"w{i}": w for i, w in enumerate(strat)},
                "median": np.median(words_left),
                "mean": np.mean(words_left),
            }
        )

    pd.DataFrame(data).to_csv("results.csv")
