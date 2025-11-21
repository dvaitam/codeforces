package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func pow2(n int) []int64 {
	res := make([]int64, n+1)
	res[0] = 1
	for i := 1; i <= n; i++ {
		res[i] = (res[i-1] * 2) % mod
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

		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		q := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &q[i])
		}

		pow := pow2(n)
		h := make([]int64, n+1)
		h[0] = 3
		for k := 1; k <= n; k++ {
			h[k] = (h[k-1] * 2) % mod
		}

		res := make([]int64, n)
		res[0] = h[p[0]] + h[q[0]] - 2
		for i := 1; i < n; i++ {
			res[i] = (res[i-1] + (h[p[i]] + h[q[i]] - 2) + mod) % mod
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
