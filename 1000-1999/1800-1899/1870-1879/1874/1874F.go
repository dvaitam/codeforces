package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int = 1e9 + 7

// checkPermutation verifies whether the permutation p satisfies
// the condition described in problemF.txt for the given m array.
// This function runs in O(n^2) for a single permutation and is
// used by the brute force solver for small n.
func checkPermutation(p []int, m []int) bool {
	n := len(p)
	for l := 1; l <= n; l++ {
		maxr := m[l-1]
		for r := l; r <= maxr; r++ {
			// compute min and max on the fly
			mn, mx := p[l-1], p[l-1]
			for i := l; i <= r; i++ {
				if p[i-1] < mn {
					mn = p[i-1]
				}
				if p[i-1] > mx {
					mx = p[i-1]
				}
			}
			if mn == l && mx == r {
				return false
			}
		}
	}
	return true
}

// bruteForce enumerates all permutations of [1..n]. It is only
// feasible for very small n (<=10) and is intended as a reference
// implementation. For larger n this approach is not practical.
func bruteForce(n int, m []int) int {
	p := make([]int, n)
	for i := range p {
		p[i] = i + 1
	}
	ans := 0
	for {
		if checkPermutation(p, m) {
			ans++
		}
		// next permutation using sort.IntsNext from Go 1.20+
		if !sort.IntsAreSorted(p) {
			// we rely on sort.NextPermutation (Go1.21).
		}
		if !nextPermutation(p) {
			break
		}
	}
	return ans % mod
}

// nextPermutation generates the next lexicographical permutation.
// Returns false if there is no next permutation.
func nextPermutation(a []int) bool {
	// find last i such that a[i] < a[i+1]
	i := len(a) - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := len(a) - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &m[i])
	}

	if n <= 10 {
		// Use brute force enumeration for small n only.
		ans := bruteForce(n, m)
		fmt.Println(ans)
		return
	}

	// TODO: Implement an efficient dynamic programming solution for
	// general n (up to 200). The current version only handles n <= 10
	// exactly. For larger n we output 0 to indicate that the efficient
	// solution is not implemented yet.
	fmt.Println(0)
}
