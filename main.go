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

func main() {
	// read in file
	content, err := ioutil.ReadFile("./words.txt")
	if err != nil {
		fmt.Println("error reading file")
		return
	}
	words := strings.Split(string(content), "\n")

	wordDist := make(map[int]int)
	for _, word := range words {
		d := distinctLetters(word)
		wordDist[d] += 1
		if d == 0 {
			fmt.Println(word)
		}
	}

	fmt.Println(wordDist)

	fmt.Println(1-jaccardSimilarity(words[:3]...), words[:3])

	// for i := 0; i < len(words); i++ {
	// 	for j := 0; j < i; j++ {
	// 		for k := 0; k < j; k++ {
	// 			// fmt.Println(words[i], words[j], words[k])
	// 		}
	// 	}
	// }

}
