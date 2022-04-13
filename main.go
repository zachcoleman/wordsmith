package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const STRAT_NUM = 3

func distinctLetters(s string) int {
	letterMap := make(map[rune]interface{})
	count := 0
	for _, letter := range s {
		if _, ok := letterMap[letter]; !ok {
			letterMap[letter] = nil
			count += 1
		}
	}
	return count
}

func jaccardSimilarity(words ...string) float64 {
	letterMap := make(map[rune]interface{})
	num, denom := 0., 0.
	for _, w := range words {
		for _, l := range w {
			if _, ok := letterMap[l]; ok {
				num += 1.
			} else {
				letterMap[l] = nil
				denom += 1.
			}

		}
	}
	return num / denom
}

func combine(n int, k int, outChan chan []int) {
	comb := []int{}
	for i := 1; i < k+1; i++ {
		comb = append(comb, i)
	}
	outChan <- comb

	for {
		didBreak := false
		i := k - 1
		for i >= 0 {
			if comb[i] != n-k+1+i {
				didBreak = true
				break
			}
			i--
		}
		if !didBreak {
			break
		}

		comb[i] += 1
		for j := i + 1; j < k; j++ {
			comb[j] = comb[j-1] + 1
		}
		tmp := make([]int, k)
		copy(tmp, comb)
		outChan <- tmp
	}
	close(outChan)
}

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
		if distinctLetters(words[i]) < 5 {
			words[i], words[j] = words[j], ""
			j--
		} else {
			i++
		}
	}
	words = words[:j+1]
	words = words[:100]

	resChan := make(chan []int)
	go combine(10, STRAT_NUM, resChan)

	for idx := range resChan {
		fmt.Println(idx)
	}

	return

	// count := 0
	// for _, indices := range combine(len(words), STRAT_NUM) {
	// 	tmpWords := make([]string, STRAT_NUM)
	// 	for arrIdx, wordIdx := range indices {
	// 		tmpWords[arrIdx] = words[wordIdx-1]
	// 	}
	// 	if jaccardSimilarity(tmpWords...) < 0.2 {
	// 		count++
	// 		fmt.Println(tmpWords)
	// 	}
	// }
	// fmt.Println(count)
}
