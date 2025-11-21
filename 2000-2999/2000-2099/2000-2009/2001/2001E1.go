package main

import (
	"bufio"
	"fmt"
	"os"
)

func powMod(base, exp, mod int64) int64 {
	result := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			result = result * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return result
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		var p int64
		fmt.Fscan(in, &n, &k, &p)

		maxHeight := n - 1
		if maxHeight < 1 {
			maxHeight = 1
		}

		F := make([][]int64, maxHeight+1)
		prefixF := make([][]int64, maxHeight+1)
		for h := 0; h <= maxHeight; h++ {
			F[h] = make([]int64, k+1)
			prefixF[h] = make([]int64, k+1)
		}

		for t := 0; t <= k; t++ {
			F[1][t] = 1 % p
			prefixF[1][t] = int64(t+1) % p
		}

		for h := 2; h <= maxHeight; h++ {
			pref := make([]int64, k+1)
			pref[0] = F[h-1][0]
			for i := 1; i <= k; i++ {
				pref[i] = (pref[i-1] + F[h-1][i]) % p
			}
			for total := 0; total <= k; total++ {
				var val int64
				for left := 0; left <= total; left++ {
					val = (val + F[h-1][left]*pref[total-left]) % p
				}
				F[h][total] = val
			}
			prefixF[h][0] = F[h][0] % p
			for i := 1; i <= k; i++ {
				prefixF[h][i] = (prefixF[h][i-1] + F[h][i]) % p
			}
		}

		dpNext := make([]int64, k+1)
		for s := 0; s <= k; s++ {
			dpNext[s] = 1 % p
		}

		for depth := n - 2; depth >= 0; depth-- {
			hSibling := n - depth - 1
			dpCurr := make([]int64, k+1)
			for total := 0; total <= k; total++ {
				var sum int64
				for child := 1; child <= total; child++ {
					maxT := total - child
					limit := child - 1
					if maxT > limit {
						maxT = limit
					}
					if maxT < 0 {
						continue
					}
					waysSibling := prefixF[hSibling][maxT]
					if waysSibling == 0 {
						continue
					}
					sum = (sum + dpNext[child]*waysSibling) % p
				}
				dpCurr[total] = sum
			}
			dpNext = dpCurr
		}

		canonical := dpNext[k] % p
		pow := powMod(2, int64(n-1), p)
		ans := canonical * pow % p
		fmt.Fprintln(out, ans)
	}
}
