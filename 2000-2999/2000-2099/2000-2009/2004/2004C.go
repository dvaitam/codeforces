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
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

		var base int64
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				base += a[i]
			} else {
				base -= a[i]
			}
		}

		var totalGap int64
		for i := 1; i < n; i += 2 {
			totalGap += a[i-1] - a[i]
		}

		if k > totalGap {
			k = totalGap
		}
		fmt.Fprintln(out, base-k)
	}
}
