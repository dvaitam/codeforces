package main

import (
	"bufio"
	"fmt"
	"os"
)

// NOTE: This is a naive reference implementation for problem D.
// It enumerates all possible substrings with exactly k ones and
// shuffles them in all possible ways to generate resulting strings.
// This approach is not efficient for the given constraints (n can
// be up to 5000) but demonstrates the logic of the operation.

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)

	if k == 0 {
		fmt.Fprintln(writer, 1)
		return
	}

	original := []byte(s)
	ones := make([]int, n+1)
	for i := 0; i < n; i++ {
		ones[i+1] = ones[i]
		if original[i] == '1' {
			ones[i+1]++
		}
	}

	seen := map[string]struct{}{}
	seen[s] = struct{}{}

	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			cnt := ones[r+1] - ones[l]
			if cnt == k {
				length := r - l + 1
				choose(length, k, func(pos []int) {
					t := append([]byte{}, original...)
					for i := l; i <= r; i++ {
						t[i] = '0'
					}
					for _, p := range pos {
						t[l+p] = '1'
					}
					seen[string(t)] = struct{}{}
				})
			}
		}
	}

	const mod = 998244353
	fmt.Fprintln(writer, len(seen)%mod)
}

// choose generates all combinations of k indices in [0, n).
func choose(n, k int, fn func([]int)) {
	comb := make([]int, k)
	var backtrack func(int, int)
	backtrack = func(start, idx int) {
		if idx == k {
			tmp := make([]int, k)
			copy(tmp, comb)
			fn(tmp)
			return
		}
		for i := start; i < n; i++ {
			comb[idx] = i
			backtrack(i+1, idx+1)
		}
	}
	backtrack(0, 0)
}
