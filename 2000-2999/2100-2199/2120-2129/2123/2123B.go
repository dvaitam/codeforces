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
		var n, j, k int
		fmt.Fscan(in, &n, &j, &k)
		a := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}
		if k == 1 {
			if a[j-1] == maxVal {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
