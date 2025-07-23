package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] > maxVal {
			maxVal = arr[i]
		}
	}

	const K = 4
	neg := -int(1e9)
	val := make([][]int, K+1)
	mod := make([][]int, K+1)
	for i := 0; i <= K; i++ {
		val[i] = make([]int, maxVal+3)
		mod[i] = make([]int, 7)
		for j := range val[i] {
			val[i][j] = neg
		}
		for j := range mod[i] {
			mod[i][j] = neg
		}
	}
	best := make([]int, K+1)
	for i := 1; i <= K; i++ {
		best[i] = neg
	}
	best[0] = 0

	for _, x := range arr {
		for k := K; k >= 1; k-- {
			cand := neg
			if x-1 >= 0 && val[k][x-1] != neg {
				if val[k][x-1]+1 > cand {
					cand = val[k][x-1] + 1
				}
			}
			if x+1 < len(val[k]) && val[k][x+1] != neg {
				if val[k][x+1]+1 > cand {
					cand = val[k][x+1] + 1
				}
			}
			if mod[k][x%7] != neg {
				if mod[k][x%7]+1 > cand {
					cand = mod[k][x%7] + 1
				}
			}
			if best[k-1] != neg {
				if best[k-1]+1 > cand {
					cand = best[k-1] + 1
				}
			}
			if cand > val[k][x] {
				val[k][x] = cand
				if cand > mod[k][x%7] {
					mod[k][x%7] = cand
				}
				if cand > best[k] {
					best[k] = cand
				}
			}
		}
	}

	if best[K] < 0 {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, best[K])
	}
}
