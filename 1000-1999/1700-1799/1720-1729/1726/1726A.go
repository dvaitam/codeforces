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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		ans := a[n-1] - a[0]
		// rotate entire array
		for i := 0; i < n-1; i++ {
			if diff := a[i] - a[i+1]; diff > ans {
				ans = diff
			}
		}
		// rotate prefix containing index 1
		minPrefix := a[0]
		for i := 0; i < n-1; i++ {
			if a[i] < minPrefix {
				minPrefix = a[i]
			}
		}
		if diff := a[n-1] - minPrefix; diff > ans {
			ans = diff
		}
		// rotate suffix containing index n
		maxSuffix := a[1]
		for i := 1; i < n; i++ {
			if a[i] > maxSuffix {
				maxSuffix = a[i]
			}
		}
		if diff := maxSuffix - a[0]; diff > ans {
			ans = diff
		}
		fmt.Fprintln(out, ans)
	}
}
