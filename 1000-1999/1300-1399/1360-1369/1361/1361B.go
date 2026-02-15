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
		if p == 1 {
			fmt.Fprintln(out, n%2)
			continue
		}
		
		sort.Slice(exps, func(i, j int) bool { return exps[i] > exps[j] })
		
		var ans int64 = 0
		var diff int64 = 0 
		var curExp int = 0
		const inf int64 = 1000005 

		for i, k := range exps {
			if i == 0 {
				ans = powMod(p, int64(k))
				diff = 1
				curExp = k
				continue
			}

			for curExp > k {
				if diff >= inf {
					curExp = k
					break
				}
				diff *= p
				curExp--
			}

			if diff > 0 {
				diff--
				term := powMod(p, int64(k))
				ans = (ans - term + mod) % mod
			} else {
				diff++
				term := powMod(p, int64(k))
				ans = (ans + term) % mod
			}
		}
		fmt.Fprintln(out, ans)
	}
}