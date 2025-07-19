package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	const INF = int64(1e9)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n+1)
		a[0] = INF
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		dp0 := make([]int64, n+1)
		dp1 := make([]int64, n+1)
		// init dp
		dp0[0] = -INF
		dp1[0] = INF
		for i := 1; i <= n; i++ {
			dp0[i] = INF
			dp1[i] = INF
		}
		// dp transition
		for i := 1; i <= n; i++ {
			// from state 0
			x := -a[i-1]
			y := dp0[i-1]
			if a[i] > x {
				dp1[i] = min(dp1[i], y)
			}
			if a[i] > y {
				dp1[i] = min(dp1[i], x)
			}
			if -a[i] > x {
				dp0[i] = min(dp0[i], y)
			}
			if -a[i] > y {
				dp0[i] = min(dp0[i], x)
			}
			// from state 1
			x = a[i-1]
			y = dp1[i-1]
			if a[i] > x {
				dp1[i] = min(dp1[i], y)
			}
			if a[i] > y {
				dp1[i] = min(dp1[i], x)
			}
			if -a[i] > x {
				dp0[i] = min(dp0[i], y)
			}
			if -a[i] > y {
				dp0[i] = min(dp0[i], x)
			}
		}
		// check
		ok := false
		now := 0
		if dp1[n] <= int64(n) {
			ok = true
			now = 1
		}
		if dp0[n] <= int64(n) {
			ok = true
			now = 0
		}
		if !ok {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		// reconstruct
		for i := n; i >= 1; i-- {
			if now == 1 {
				x := -a[i-1]
				y := dp0[i-1]
				if a[i] > x && dp1[i] == y {
					now = 0
					continue
				}
				if a[i] > y && dp1[i] == x {
					now = 0
					continue
				}
				x = a[i-1]
				y = dp1[i-1]
				if a[i] > x && dp1[i] == y {
					now = 1
					continue
				}
				if a[i] > y && dp1[i] == x {
					now = 1
					continue
				}
			} else {
				x := -a[i-1]
				y := dp0[i-1]
				if -a[i] > x && dp0[i] == y {
					now = 0
				}
				if -a[i] > y && dp0[i] == x {
					now = 0
				}
				x = a[i-1]
				y = dp1[i-1]
				if -a[i] > x && dp0[i] == y {
					now = 1
				}
				if -a[i] > y && dp0[i] == x {
					now = 1
				}
				a[i] = -a[i]
			}
		}
		// output modified a[1..n]
		for i := 1; i <= n; i++ {
			fmt.Fprintf(writer, "%d ", a[i])
		}
		fmt.Fprintln(writer)
	}
}
