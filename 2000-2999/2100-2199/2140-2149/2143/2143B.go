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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		b := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &b[i])
		}

		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

		// First, assign one most expensive item to each voucher.
		ans := int64(0)
		ptr := 0
		for i := 0; i < k; i++ {
			ans += a[i] // pay for the most expensive in the group
			if b[i] == 1 {
				// Free item is itself
				ans += a[i]
			} else {
				// Will add the free item later
			}
		}

		ptr = k
		for i := k - 1; i >= 0; i-- {
			freePos := ptr + b[i] - 2
			ans += a[freePos]
			ptr = freePos + 1
		}

		fmt.Fprintln(out, ans)
	}
}
