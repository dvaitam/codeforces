package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, k, q int
		fmt.Fscan(in, &n, &m, &k, &q)
		xs := make([]int, q)
		ys := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &xs[i], &ys[i])
		}

		rowUsed := make([]bool, n+1)
		colUsed := make([]bool, m+1)
		cntRow, cntCol := 0, 0
		ans := int64(1)
		for i := q - 1; i >= 0; i-- {
			if cntRow == n || cntCol == m {
				break
			}
			x, y := xs[i], ys[i]
			if !rowUsed[x] || !colUsed[y] {
				ans = ans * int64(k) % mod
				if !rowUsed[x] {
					rowUsed[x] = true
					cntRow++
				}
				if !colUsed[y] {
					colUsed[y] = true
					cntCol++
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
