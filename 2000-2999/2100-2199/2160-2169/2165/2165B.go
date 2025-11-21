package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
		}

		counts := make([]int, 0, n)
		for _, c := range freq {
			if c > 0 {
				counts = append(counts, c)
			}
		}

		sort.Ints(counts)
		m := len(counts)

		prod := int64(1)
		inv := make([]int64, m)
		total := 0
		for i, c := range counts {
			prod = prod * int64(c) % mod
			inv[i] = modPow(int64(c), mod-2)
			total += c
		}

		dp := make([]int64, total+1)
		dp[0] = 1
		ans := int64(1) // empty subset T

		for i, c := range counts {
			W := total - 2*c
			if W >= 0 {
				prefix := int64(0)
				for s := 0; s <= W; s++ {
					prefix += dp[s]
					if prefix >= mod {
						prefix -= mod
					}
				}
				ans = (ans + inv[i]*prefix) % mod
			}

			invc := inv[i]
			for s := total - c; s >= 0; s-- {
				if dp[s] == 0 {
					continue
				}
				val := dp[s] * invc % mod
				dp[s+c] += val
				if dp[s+c] >= mod {
					dp[s+c] -= mod
				}
			}
		}

		ans = ans * prod % mod
		fmt.Fprintln(out, ans)
	}
}
