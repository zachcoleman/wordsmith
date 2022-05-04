# Wordsmith

## Problem

The popular game Wordle (New York Times) is a word puzzle designed to find a hidden word. An optimal starting strategy is a set of words such that the number of possible hidden words is minimized.

## Strategy Space

The word bank in `data/` contains 15,918 5-letter words. If a 3 word starting strategy is considered the possible strategies are C(15,918, 3) ~= 672B.

The python package (wordsmith) provides code to evaluate a strategy by comparing the distribution of remaining possibilities over a sample of hidden words for a strategy (i.e.how well a set of guesses narrows down the possible hidden words).

This empirically driven evaluation method makes it infeasible to reasonably calculate/parse the entire ~672B combinations. The following assumption is made to reduce the search space: an optimal strategy will contain only unique letters across the words.

To even evaluate the letter-uniqueness of a set of words can be quite cumbersome to compute for this many possible combinations. Therefore, Go is used for its fast computing and multithreading to parse through a sample of combinations. 

The Go program `candidate_filtering` parses the approximately 10,000 words with a set of distinct letters down to a sample of 1,000 words in which all combinations of 3 can be chosen (C(1,000, 3) ~= 116M). These combinations are filtered down to sets with a Jaccard similarity of 0 resulting in finally a more reasonable 215K strategies to be evaluated.

## Evaluation

## Results
