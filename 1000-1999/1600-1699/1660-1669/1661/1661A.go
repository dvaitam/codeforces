package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt.
// For each index we may swap a[i] and b[i]. The optimal strategy
// is to sort each pair so that a[i] <= b[i], which minimizes
// |a[i]-a[i+1]| + |b[i]-b[i+1]| for all adjacent pairs.
// After sorting, we simply sum these absolute differences.

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
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
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		for i := 0; i < n; i++ {
			if a[i] > b[i] {
				a[i], b[i] = b[i], a[i]
			}
		}
		var ans int64
		for i := 0; i < n-1; i++ {
			ans += abs64(a[i]-a[i+1]) + abs64(b[i]-b[i+1])
		}
		fmt.Fprintln(out, ans)
	}
}
