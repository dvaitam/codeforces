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
		a := make([]int64, n)
		var maxVal int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}
		var ops int64
		if n > 0 {
			ops = a[0]
		}
		for i := 1; i < n; i++ {
			if a[i] > a[i-1] {
				ops += a[i] - a[i-1]
			}
		}
		if ops == maxVal {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
