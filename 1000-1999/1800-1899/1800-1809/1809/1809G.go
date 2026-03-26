package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	fmt.Fscan(reader, &n, &k)

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	if n == 1 {
		fmt.Println(1)
		return
	}

	// Brute force: enumerate all permutations, check if all games are predictable.
	// A game is predictable iff |champion - challenger| > k.
	// The champion is always the max-rated among those processed so far
	// (since when |diff|>k, the higher-rated player wins).
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i
	}
	count := int64(0)
	var generate func(int)
	generate = func(start int) {
		if start == n {
			predictable := true
			curMax := a[perm[0]]
			for i := 1; i < n; i++ {
				challenger := a[perm[i]]
				diff := curMax - challenger
				if diff < 0 {
					diff = -diff
				}
				if diff <= k {
					predictable = false
					break
				}
				if challenger > curMax {
					curMax = challenger
				}
			}
			if predictable {
				count++
			}
			return
		}
		for i := start; i < n; i++ {
			perm[start], perm[i] = perm[i], perm[start]
			generate(start + 1)
			perm[start], perm[i] = perm[i], perm[start]
		}
	}
	generate(0)
	fmt.Println(count % MOD)
}
