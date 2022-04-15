package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"wordsmith/go-utils"
)

const STRAT_NUM = 3

func main() {
	// read in file of words
	content, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		fmt.Println("error reading file")
		return
	}
	words := strings.Split(string(content), "\n")

	// filter to only unique words
	i, j := 0, len(words)-1
	for i < j {
		if utils.DistinctLetters(words[i]) < 5 {
			words[i], words[j] = words[j], ""
			j--
		} else {
			i++
		}
	}
	words = words[:j+1]
	words = words[:5]

	num_workers := 2
	inQueue := make(chan []string, num_workers*2)
	outQueue := make(chan []string, num_workers*2)
	var wg sync.WaitGroup

	for i := 0; i < num_workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tmpWords := range inQueue {
				if utils.JaccardSimilarity(tmpWords...) < 0.2 {
					outQueue <- tmpWords
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(outQueue)
	}()

	total, left := 0, 0
	go func() {
		for indices := range utils.IndexCombinations(len(words), STRAT_NUM) {
			fmt.Println(indices)
			tmpWords := make([]string, STRAT_NUM)
			for arrIdx, wordIdx := range indices {
				tmpWords[arrIdx] = words[wordIdx-1]
			}
			total++
			inQueue <- tmpWords
		}
		close(inQueue)
	}()

	for range outQueue {
		left++
	}

	fmt.Println(total, left)
}
