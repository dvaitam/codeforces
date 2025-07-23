package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)

	dpPrev := make([][]int64, n+1)
	dpCurr := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		dpPrev[i] = make([]int64, k+1)
		dpCurr[i] = make([]int64, k+1)
	}
	dpPrev[0][0] = 1

	for i := 0; i < n; i++ {
		diff := 0
		if i > 0 {
			diff = a[i] - a[i-1]
		}
		for j := 0; j <= n; j++ {
			for s := 0; s <= k; s++ {
				dpCurr[j][s] = 0
			}
		}
		for open := 0; open <= i; open++ {
			for s := 0; s <= k; s++ {
				val := dpPrev[open][s]
				if val == 0 {
					continue
				}
				newS := s + open*diff
				if newS > k {
					continue
				}
				dpCurr[open][newS] = (dpCurr[open][newS] + val) % MOD
				dpCurr[open+1][newS] = (dpCurr[open+1][newS] + val) % MOD
				if open > 0 {
					mult := val * int64(open) % MOD
					dpCurr[open][newS] = (dpCurr[open][newS] + mult) % MOD
					dpCurr[open-1][newS] = (dpCurr[open-1][newS] + mult) % MOD
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}

	ans := int64(0)
	for s := 0; s <= k; s++ {
		ans = (ans + dpPrev[0][s]) % MOD
	}
	fmt.Fprintln(writer, ans)
}
