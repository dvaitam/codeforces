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
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		d := make([]int64, n)
		for i := 1; i < n; i++ {
			d[i] = a[i+1] - a[i]
		}
		// weights for moving to the right or left between adjacent cities
		wR := make([]int64, n)   // weight from i to i+1 for i in [1,n-1]
		wL := make([]int64, n+1) // weight from i to i-1 for i in [2,n]
		for i := 1; i < n; i++ {
			wR[i] = d[i]
			wL[i+1] = d[i]
		}
		// closest city edges
		wR[1] = 1
		wL[n] = 1
		for i := 2; i <= n-1; i++ {
			left := a[i] - a[i-1]
			right := a[i+1] - a[i]
			if left < right {
				wL[i] = 1
			} else {
				wR[i] = 1
			}
		}
		// prefix sums for fast queries
		prefR := make([]int64, n)
		for i := 1; i < n; i++ {
			prefR[i] = prefR[i-1] + wR[i]
		}
		prefL := make([]int64, n+1)
		for i := 2; i <= n; i++ {
			prefL[i] = prefL[i-1] + wL[i]
		}

		var m int
		fmt.Fscan(in, &m)
		for ; m > 0; m-- {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x < y {
				fmt.Fprintln(out, prefR[y-1]-prefR[x-1])
			} else {
				fmt.Fprintln(out, prefL[x]-prefL[y])
			}
		}
	}
}
