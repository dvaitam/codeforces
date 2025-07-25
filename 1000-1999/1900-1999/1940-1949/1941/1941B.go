package main

import (
	"bufio"
	"fmt"
	"os"
)

// canZero checks if array a can be reduced to all zeros using the allowed operation.
func canZero(a []int64) bool {
	n := len(a)
	for i := 0; i <= n-3; i++ {
		x := a[i]
		if a[i+1] < 2*x || a[i+2] < x {
			return false
		}
		a[i+1] -= 2 * x
		a[i+2] -= x
		// a[i] becomes zero implicitly
	}
	return a[n-2] == 0 && a[n-1] == 0
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
		if canZero(a) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
