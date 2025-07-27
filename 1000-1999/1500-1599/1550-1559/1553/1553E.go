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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			p[i]-- // zero-index values
		}

		cnt := make([]int, n)
		for i := 0; i < n; i++ {
			// compute shift k that puts p[i] at position i after shift
			k := (i - p[i] + n) % n
			cnt[k]++
		}

		var candidates []int
		for k := 0; k < n; k++ {
			if n-cnt[k] <= 2*m { // mismatches can be fixed with at most m swaps
				if checkShift(p, n, m, k) {
					candidates = append(candidates, k)
				}
			}
		}

		fmt.Fprint(out, len(candidates))
		for _, k := range candidates {
			fmt.Fprintf(out, " %d", k)
		}
		fmt.Fprintln(out)
	}
}

func checkShift(p []int, n, m, k int) bool {
	// permutation mapping from index to where element should go
	visited := make([]bool, n)
	swaps := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			j := i
			cycleLen := 0
			for !visited[j] {
				visited[j] = true
				j = (p[j] + k) % n
				cycleLen++
			}
			swaps += cycleLen - 1
			if swaps > m {
				return false
			}
		}
	}
	return swaps <= m
}
