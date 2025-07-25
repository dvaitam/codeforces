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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		prefix := int64(0)
		bestB := 0
		ans := int64(0)
		limit := k
		if n < limit {
			limit = n
		}
		for i := 0; i < limit; i++ {
			prefix += int64(a[i])
			if b[i] > bestB {
				bestB = b[i]
			}
			remaining := k - (i + 1)
			if remaining < 0 {
				remaining = 0
			}
			total := prefix + int64(remaining*bestB)
			if total > ans {
				ans = total
			}
		}
		fmt.Fprintln(out, ans)
	}
}
