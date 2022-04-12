package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

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

func combine(n int, k int) [][]int {

	if k == 1 {
		ret := [][]int{}
		for n > 0 {
			ret = append(ret, []int{n})
			n--
		}
		return ret
	}

	ret := [][]int{}
	for n > 1 {
		set := combine(n-1, k-1)
		for _, subset := range set {
			tmp := append([]int{n}, subset...)
			ret = append(ret, tmp)
		}
		n--
	}

	return ret
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

	const STRAT_NUM = 3
	count := 0
	for _, indices := range combine(len(words), STRAT_NUM) {
		tmpWords := make([]string, STRAT_NUM)
		for arrIdx, wordIdx := range indices {
			tmpWords[arrIdx] = words[wordIdx-1]
		}
		if jaccardSimilarity(tmpWords...) < 0.2 {
			count++
			fmt.Println(tmpWords)
		}
	}

	fmt.Println(count)
}
