package utils

import (
	"fmt"
	"sync"
)

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

func Factorial(n int) int {
	if n < 0 {
		return -1
	} else if n <= 1 {
		return 1
	}
	ret := 1
	for i := 1; i <= n; i++ {
		ret *= i
	}
	return ret
}

func Choose(n int, k int) int {
	num, denom := 1, 1

	// Example: C(10, 3)
	// since 7 > 3 want to eliminate
	// 7 multiplications instead of 3:
	// this:     10*9*8 / 3*2*1
	// not this: 10*9*8*7*6*5*4 / 7*6*5*4*3*2*1
	//
	// for C(10, 7)
	if n-k > k {
		for i := n; i > n-k; i-- {
			num *= i
		}
		for i := 1; i <= k; i++ {
			denom *= i
		}
	} else {
		for i := n; i > k; i-- {
			num *= i
		}
		for i := 1; i <= n-k; i++ {
			denom *= i
		}
	}
	return num / denom
}

func IndexCombinations(n int, k int) chan []int {
	incrementCounts, rectifyCounts := 0, 0
	ret := make(chan []int, 2*k+1)
	go func() {
		comb := []int{}
		for i := 1; i < k+1; i++ {
			comb = append(comb, i)
		}
		tmp := make([]int, k)
		copy(tmp, comb)
		ret <- tmp

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
			/*
				Loop rectifies the right-side digits
				after a carryover increment e.g.
					[1, 2, 7] -> [1, 3, 7]
				but in reality we need to go to:
					[1, 3, 4]
			*/

			if i+1 < k {
				rectifyCounts++
			} else {
				incrementCounts++
			}

			// This seems to be most expensive operation and
			// largest bottleneck for generating combinations
			// and the problem is at worst for C(n, n/2)
			// (the maintenance of lexicographical order may be to blame
			// in the iterative generation process)
			// Note: rectifyCounts vs incrementCounts move inversely
			// as k -> 0, incrementCounts increase and rectifyCounts decrease
			// as k -> n, incrementCounts decrease and rectifyCounts increase
			for j := i + 1; j < k; j++ {
				comb[j] = comb[j-1] + 1
			}

			tmp := make([]int, k)
			copy(tmp, comb)
			ret <- tmp
		}
		close(ret)
		fmt.Println("rectifyCounts: ", rectifyCounts)
		fmt.Println("incrementCounts: ", incrementCounts)
	}()

	return ret
}

func MultiIndexCombinations(n int, k int) chan []int {
	ret := make(chan []int, 2*k+1)
	var wg sync.WaitGroup
	for firstDigit := 1; firstDigit <= n-k+1; firstDigit++ {
		wg.Add(1)
		go func(firstDigit int) {
			comb := []int{}
			for i := 0; i < k; i++ {
				comb = append(comb, firstDigit+i)
			}
			tmp := make([]int, k)
			copy(tmp, comb)
			ret <- tmp

			for {
				didBreak := false
				i := k - 1
				for i > 0 {
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
			wg.Done()
		}(firstDigit)
	}
	go func() {
		wg.Wait()
		close(ret)
	}()

	return ret
}

func splitDomain(n int, k int, numWorkers int) [][2]int {
	quants := []float32{}
	for i := 1; i <= numWorkers; i++ {
		quants = append(quants, float32(i)/float32(numWorkers))
	}

	qIdx, prev := 0, 0
	breakPts := [][2]int{}

	runningTotal := 0
	total := Choose(n, k)
	for i := 1; i <= n-k+1; i++ {
		runningTotal += Choose(n-i, k-1)
		if qIdx < len(quants) && float32(runningTotal)/float32(total) >= quants[qIdx] {
			fmt.Println(quants[qIdx], i, float32(runningTotal)/float32(total))
			breakPts = append(breakPts, [2]int{prev, i})
			prev = i
			qIdx++
		}
	}

	return breakPts
}

func IndexCombinationsMulti(n int, k int, numWorkers int) chan []int {

	// if not possible to divide up leftmost digits
	// into numWorkers sets to work then revise numWorkers
	// this implies this strategy of dividing up combo generation
	// only works effectively when k is sufficiently less than n (k < n/2)
	// in case when k = n/2 we suboptimally are bottlenecked at 1 thread
	// responsible for 50% of work (this is worst case scenario)
	// TODO: implement strategy for arbitrary combination of n, k and/or
	// k closer to n like C(10, 7) where 70% of combinations have 1 as
	// leftmost digit
	// Alternative solution: is to form sets from complement problem
	// C(10, 7) is equivalent C(10, 3) s.t. instead of choosing 7 to
	// include in set 3 are being chosen to be excluded this
	// reframing of problem is more memory efficient if put back on caller
	// or convenience method could be made s.t the complement is automatically
	// calculated and yielded back

	breakPts := splitDomain(n, k, numWorkers)
	if len(breakPts) < numWorkers {
		numWorkers = len(breakPts)
	}

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
