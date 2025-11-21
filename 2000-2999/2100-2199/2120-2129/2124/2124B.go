package main

import (
	"bufio"
	"fmt"
	"os"
)

func min64(a, b int64) int64 {
	if a < b {
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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		pre := make([]int64, n+1)
		curMin := int64(1 << 60)
		for i := 0; i < n; i++ {
			if a[i] < curMin {
				curMin = a[i]
			}
			pre[i+1] = pre[i] + curMin
		}

		ans := pre[n]
		if n >= 2 {
			ans = min64(ans, a[0]+a[1])
		}
		for pos := 2; pos < n; pos++ {
			ans = min64(ans, pre[pos])
		}
		fmt.Fprintln(out, ans)
	}
}
