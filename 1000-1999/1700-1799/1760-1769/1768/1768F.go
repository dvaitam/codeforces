package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}

	// build sparse table for range minimum queries
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	st[0] = make([]int, n)
	copy(st[0], a)
	for k := 1; k < K; k++ {
		size := n - (1 << k) + 1
		st[k] = make([]int, size)
		for i := 0; i < size; i++ {
			x := st[k-1][i]
			y := st[k-1][i+(1<<(k-1))]
			if x < y {
				st[k][i] = x
			} else {
				st[k][i] = y
			}
		}
	}
	query := func(l, r int) int {
		if l > r {
			l, r = r, l
		}
		k := log[r-l+1]
		x := st[k][l]
		y := st[k][r-(1<<k)+1]
		if x < y {
			return x
		}
		return y
	}

	dp := make([]int64, n)
	const INF int64 = 1 << 60
	for i := 1; i < n; i++ {
		dp[i] = INF
	}

	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			m := query(j, i)
			d := int64(i - j)
			cost := int64(m) * d * d
			if dp[j]+cost < dp[i] {
				dp[i] = dp[j] + cost
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range dp {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
