package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, k int
	fmt.Fscan(in, &n, &k)
	a := make([]int, n)
	b := make([]int, n)
	sumR := 0
	sumB := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i], &b[i])
		sumR += a[i]
		sumB += b[i]
	}
	dp := make([]bool, k)
	dp[0] = true
	for i := 0; i < n; i++ {
		ndp := make([]bool, k)
		for r := 0; r < k; r++ {
			if !dp[r] {
				continue
			}
			maxx := min(a[i], k-1)
			for x := 0; x <= maxx; x++ {
				y := (k - ((r + x) % k)) % k
				if y <= b[i] && x+y <= a[i]+b[i] {
					ndp[(r+x)%k] = true
				}
			}
		}
		dp = ndp
	}
	total := sumR + sumB
	res := total / k
	if !dp[sumR%k] {
		res--
	}
	if res < 0 {
		res = 0
	}
	fmt.Fprintln(out, res)
}
