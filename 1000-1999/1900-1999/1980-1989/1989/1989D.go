package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 1000005
const INF int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	d := make([]int64, N)
	dp := make([]int64, N)
	for i := range d {
		d[i] = INF
	}

	var n, m int
	fmt.Fscan(in, &n, &m)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		var b int
		fmt.Fscan(in, &b)
		diff := int64(a[i] - b)
		if diff < d[a[i]] {
			d[a[i]] = diff
		}
	}
	for i := 1; i < N; i++ {
		if d[i-1] < d[i] {
			d[i] = d[i-1]
		}
		if d[i] <= int64(i) {
			dp[i] = 2 + dp[i-int(d[i])]
		} else {
			dp[i] = 0
		}
	}
	var ans int64
	for ; m > 0; m-- {
		var c int64
		fmt.Fscan(in, &c)
		if c >= N {
			x := (c-N)/d[N-1] + 1
			c -= x * d[N-1]
			ans += 2 * x
		}
		ans += dp[int(c)]
	}
	fmt.Fprint(out, ans)
}
