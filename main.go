package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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

	// shuffle and sample data
	for i := 0; i < len(words); i++ {
		j, k := rand.Intn(len(words)), rand.Intn(len(words))
		words[j], words[k] = words[k], words[j]
	}
	words = words[:1_000]

	// define threading params
	numWorkers := 4
	inQueue := make(chan []string, numWorkers*10)
	outQueue := make(chan []string, numWorkers*10)
	var evalWg sync.WaitGroup
	var indexWg sync.WaitGroup

	// generate indices and define word combinations
	total, left := 0, 0
	combChan := utils.IndexCombinationsMulti(len(words), STRAT_NUM, numWorkers)
	for i := 0; i < numWorkers; i++ {
		indexWg.Add(1)
		go func() {
			defer indexWg.Done()
			for indices := range combChan {
				tmpWords := make([]string, STRAT_NUM)
				for arrIdx, wordIdx := range indices {
					tmpWords[arrIdx] = words[wordIdx-1]
				}
				total++
				inQueue <- tmpWords
			}
		}()
	}
	go func() {
		indexWg.Wait()
		close(inQueue)
	}()

	// start evaluation workers
	for i := 0; i < numWorkers; i++ {
		evalWg.Add(1)
		go func() {
			defer evalWg.Done()
			for tmpWords := range inQueue {
				if utils.JaccardSimilarity(tmpWords...) == 0 {
					// fmt.Println(tmpWords)
					outQueue <- tmpWords
				}
			}
		}()
	}
	go func() {
		evalWg.Wait()
		close(outQueue)
	}()

	// process outputs
	f, err := os.Create("candidates.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for candWords := range outQueue {
		left++
		f.WriteString(strings.Join(candWords, ",") + "\n")
	}
	fmt.Println(total, "->", left)
}
