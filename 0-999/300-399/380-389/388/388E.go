package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}

	const inf int64 = -1 << 60
	dp := [2]int64{0, inf}
	for i := 0; i < n-1; i++ {
		next := [2]int64{inf, inf}
		for px := 0; px <= 1; px++ {
			if dp[px] <= inf/2 {
				continue
			}
			for xi := 0; xi <= 1; xi++ {
				count := px + (1 - xi)
				var joy int64
				if count == 0 {
					joy = a[i]
				} else if count == 1 {
					joy = b[i]
				} else {
					joy = c[i]
				}
				val := dp[px] + joy
				if val > next[xi] {
					next[xi] = val
				}
			}
		}
		dp = next
	}

	ans := dp[0] + a[n-1]
	if v := dp[1] + b[n-1]; v > ans {
		ans = v
	}
	fmt.Fprintln(out, ans)
}
