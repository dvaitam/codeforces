package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for contest 1628A.
// It constructs the lexicographically maximum array b by repeatedly removing
// a prefix whose MEX becomes the next element of b. For each step we compute
// the MEX of the remaining array and take the smallest prefix that contains
// all numbers from 0..mex-1. This greedy approach yields the optimal result.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// Count occurrences of every value.
		freq := make([]int, n+2)
		for _, v := range a {
			if v <= n+1 {
				freq[v]++
			}
		}

		res := make([]int, 0)
		for i := 0; i < n; {
			// Compute current MEX of the remaining array.
			mex := 0
			for ; mex <= n+1 && freq[mex] > 0; mex++ {
			}

			if mex == 0 {
				// Remove one element and output 0.
				freq[a[i]]--
				res = append(res, 0)
				i++
				continue
			}

			// Need all numbers 0..mex-1 in the next segment.
			seen := make([]bool, mex)
			need := mex
			j := i
			for j < n && need > 0 {
				v := a[j]
				freq[v]--
				if v < mex && !seen[v] {
					seen[v] = true
					need--
				}
				j++
			}
			res = append(res, mex)
			i = j
		}

		fmt.Fprintln(out, len(res))
		for idx, v := range res {
			if idx > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
