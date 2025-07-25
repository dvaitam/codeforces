package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var L, R, K int
	if _, err := fmt.Fscan(in, &n, &L, &R, &K); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	minVal := make([]int, n)
	maxVal := make([]int, n)
	// allowed range for last element
	minVal[n-1] = L
	maxVal[n-1] = R
	// compute feasible ranges backwards
	for i := n - 2; i >= 0; i-- {
		var lo, hi int
		if a[i] < a[i+1] {
			lo = L
			if minVal[i+1]-K > lo {
				lo = minVal[i+1] - K
			}
			hi = R
			if maxVal[i+1]-1 < hi {
				hi = maxVal[i+1] - 1
			}
		} else if a[i] > a[i+1] {
			lo = L
			if minVal[i+1]+1 > lo {
				lo = minVal[i+1] + 1
			}
			hi = R
			if maxVal[i+1]+K < hi {
				hi = maxVal[i+1] + K
			}
		} else {
			lo = L
			if minVal[i+1] > lo {
				lo = minVal[i+1]
			}
			hi = R
			if maxVal[i+1] < hi {
				hi = maxVal[i+1]
			}
		}
		if lo > hi {
			fmt.Println(-1)
			return
		}
		minVal[i], maxVal[i] = lo, hi
	}

	// reconstruct lexicographically smallest sequence
	b := make([]int, n)
	b[0] = minVal[0]
	for i := 0; i < n-1; i++ {
		if a[i] < a[i+1] {
			lo := b[i] + 1
			if minVal[i+1] > lo {
				lo = minVal[i+1]
			}
			hi := b[i] + K
			if maxVal[i+1] < hi {
				hi = maxVal[i+1]
			}
			if lo > hi {
				fmt.Println(-1)
				return
			}
			b[i+1] = lo
		} else if a[i] > a[i+1] {
			lo := b[i] - K
			if minVal[i+1] > lo {
				lo = minVal[i+1]
			}
			hi := b[i] - 1
			if maxVal[i+1] < hi {
				hi = maxVal[i+1]
			}
			if lo > hi {
				fmt.Println(-1)
				return
			}
			b[i+1] = lo
		} else {
			if b[i] < minVal[i+1] || b[i] > maxVal[i+1] {
				fmt.Println(-1)
				return
			}
			b[i+1] = b[i]
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, b[i])
	}
	out.Flush()
}
