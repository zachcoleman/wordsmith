package utils

import "sync"

func DistinctLetters(s string) int {
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

func JaccardSimilarity(words ...string) float64 {
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

func IndexCombinations(n int, k int) chan []int {
	ret := make(chan []int, 8) // small performance lift
	go func() {
		comb := []int{}
		for i := 1; i < k+1; i++ {
			comb = append(comb, i)
		}
		ret <- comb

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
			ret <- tmp
		}
		close(ret)
	}()
	return ret
}

func IndexCombinationsMulti(n int, k int, numWorkers int) chan []int {
	ret := make(chan []int, numWorkers*4)
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			comb := []int{}
			for i := 1; i < k+1; i++ {
				comb = append(comb, i)
			}
			ret <- comb

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
				ret <- tmp
			}
			close(ret)
		}()
	}

	return ret
}
