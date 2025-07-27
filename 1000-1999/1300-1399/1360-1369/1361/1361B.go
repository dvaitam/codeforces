package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1000000007

func powMod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
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
		var p int64
		fmt.Fscan(in, &n, &p)
		exps := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &exps[i])
		}
		sort.Slice(exps, func(i, j int) bool { return exps[i] > exps[j] })
		if p == 1 {
			if n%2 == 0 {
				fmt.Fprintln(out, 0)
			} else {
				fmt.Fprintln(out, 1)
			}
			continue
		}
		const inf = 1000000
		diff := 0
		curExp := 0
		var ans int64
		for i, k := range exps {
			if diff == 0 {
				diff = 1
				curExp = k
				ans = powMod(p, int64(k))
				continue
			}
			for diff < inf && curExp > k {
				diff *= int(p)
				if diff >= inf {
					break
				}
				ans = ans * p % mod
				curExp--
			}
			if diff >= inf {
				// Difference is already huge, just subtract remaining contributions
				for j := i; j < n; j++ {
					ans = (ans - powMod(p, int64(exps[j]))%mod + mod) % mod
				}
				diff = 0
				break
			}
			diff--
			ans = (ans - powMod(p, int64(k))%mod + mod) % mod
			curExp = k
		}
		fmt.Fprintln(out, (ans%mod+mod)%mod)
	}
}
