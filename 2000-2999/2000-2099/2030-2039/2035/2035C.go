package main

import (
	"bufio"
	"fmt"
	"math/bits"
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
		var n int
		fmt.Fscan(in, &n)
		ans, perm := buildPermutation(n)
		fmt.Fprintln(out, ans)
		for i, v := range perm {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

// buildPermutation returns the optimal value of k and one permutation that attains it.
func buildPermutation(n int) (int, []int) {
	switch n {
	case 1:
		return 0, []int{1}
	case 2:
		return 2, []int{1, 2}
	case 3:
		return 2, []int{1, 2, 3}
	case 4:
		return 6, []int{1, 2, 3, 4}
	}

	if n%2 == 1 {
		return n, oddPermutation(n)
	}

	// Even n.
	if n&(n-1) == 0 { // power of two
		ans := (1 << bits.Len(uint(n))) - 1
		return ans, oddPermutation(n)
	}

	highestPower := 1 << (bits.Len(uint(n)) - 1)
	target := (highestPower << 1) - 1 // all bits up to the MSB set
	missing := target &^ n            // bits not present in n that we still need

	penultimate := highestPower - 1 // has every lower bit set
	last := n

	perm := make([]int, 0, n)
	used := map[int]bool{
		missing:     true,
		penultimate: true,
		last:        true,
	}

	for i := 1; i <= n; i++ {
		if !used[i] {
			perm = append(perm, i)
		}
	}

	perm = append(perm, missing, penultimate, last)
	return target, perm
}

// oddPermutation produces [2, 1, 3, 4, ..., n], which is optimal for all odd n
// (and also works for several even cases).
func oddPermutation(n int) []int {
	perm := make([]int, 0, n)
	if n >= 2 {
		perm = append(perm, 2, 1)
	}
	for i := 3; i <= n; i++ {
		perm = append(perm, i)
	}
	return perm
}
