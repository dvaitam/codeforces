package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

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
		used := make([]bool, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > 0 {
				used[a[i]] = true
			}
		}

		// Fill zeros greedily: put missing numbers in increasing order into zeros.
		missing := make([]int, 0)
		for v := 1; v <= n; v++ {
			if !used[v] {
				missing = append(missing, v)
			}
		}
		sort.Ints(missing)
		idxMissing := 0
		for i := 0; i < n; i++ {
			if a[i] == 0 {
				a[i] = missing[idxMissing]
				idxMissing++
			}
		}

		// Compute cost of resulting permutation
		sorted := make([]int, n)
		copy(sorted, a)
		sort.Ints(sorted)

		cost := 0
		if !equal(a, sorted) {
		L1:
			for length := 1; length <= n; length++ {
				for l := 0; l+length <= n; l++ {
					r := l + length
					segment := append([]int{}, a[l:r]...)
					sort.Ints(segment)
					newArr := append(append(append([]int{}, a[:l]...), segment...), a[r:]...)
					if equal(newArr, sorted) {
						cost = length
						break L1
					}
				}
			}
		}

		fmt.Fprintln(out, cost)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
