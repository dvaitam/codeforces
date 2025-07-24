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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		pre := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pre[i] = pre[i-1]
			if a[i] < i {
				pre[i]++
			}
		}
		var ans int64
		for j := 1; j <= n; j++ {
			if a[j] < j {
				idx := a[j] - 1
				if idx >= 0 {
					if idx > n {
						idx = n
					}
					ans += int64(pre[idx])
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
