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

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)

		a := make([]int64, n)
		for i := int64(0); i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// Build the frequency/value map
		mp := make(map[int64]int64, n)
		for i := int64(0); i < n; i++ {
			var v int64
			fmt.Fscan(in, &v)
			mp[a[i]] = v
		}

		var mx int64
		for i := int64(0); i < n; i++ {
			var ans int64

			f := min(m/a[i], mp[a[i]])
			f1 := min((m-f*a[i])/(a[i]+1), mp[a[i]+1])
			ans = f*a[i] + f1*(a[i]+1)

			f3 := min(f, mp[a[i]+1]-f1)
			ans += f3

			mx = max(mx, min(ans, m))
			if mx == m {
				break
			}
		}
		fmt.Fprintln(out, mx)
	}
}
