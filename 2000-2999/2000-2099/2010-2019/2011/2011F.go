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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		diffs := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			diffs[i] = abs(a[i+1] - a[i])
		}
		total := int64(n)
		var cur int64 = 0
		for i := 0; i < len(diffs); i++ {
			if diffs[i] == 1 {
				cur++
			} else {
				cur = 0
			}
			total += cur
		}
		fmt.Fprintln(out, total)
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
