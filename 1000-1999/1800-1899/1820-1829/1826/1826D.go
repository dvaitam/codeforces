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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		b := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &b[i])
		}
		pre := make([]int64, n+2)
		const negInf int64 = -1 << 60
		pre[0] = negInf
		for i := 1; i <= n; i++ {
			val := b[i] + int64(i)
			if val > pre[i-1] {
				pre[i] = val
			} else {
				pre[i] = pre[i-1]
			}
		}
		suf := make([]int64, n+2)
		suf[n+1] = negInf
		for i := n; i >= 1; i-- {
			val := b[i] - int64(i)
			if val > suf[i+1] {
				suf[i] = val
			} else {
				suf[i] = suf[i+1]
			}
		}
		var ans int64 = negInf
		for j := 2; j <= n-1; j++ {
			val := pre[j-1] + b[j] + suf[j+1]
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(out, ans)
	}
}
