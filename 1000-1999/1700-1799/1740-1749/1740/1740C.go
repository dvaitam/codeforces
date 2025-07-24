package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		var ans int64
		for i := 0; i <= n-3; i++ {
			diff := (a[n-1] - a[i]) + (a[i+1] - a[i])
			if diff > ans {
				ans = diff
			}
		}
		for i := 2; i < n; i++ {
			diff := (a[i] - a[0]) + (a[i] - a[i-1])
			if diff > ans {
				ans = diff
			}
		}
		fmt.Fprintln(out, ans)
	}
}
