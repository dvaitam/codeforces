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

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}

	for ; q > 0; q-- {
		var typ, k int
		var d int64
		fmt.Fscan(in, &typ, &k, &d)
		if typ == 1 {
			for i := 0; i < k; i++ {
				idx := n - k + i
				a[idx] += d * int64(i+1)
			}
		} else {
			for i := 0; i < k; i++ {
				idx := m - k + i
				b[idx] += d * int64(i+1)
			}
		}

		// compute minimal path using greedy for convex arrays
		i, j := 0, 0
		res := a[0] + b[0]
		for i < n-1 && j < m-1 {
			da := a[i+1] - a[i]
			db := b[j+1] - b[j]
			if da < db {
				i++
			} else {
				j++
			}
			res += a[i] + b[j]
		}
		for i < n-1 {
			i++
			res += a[i] + b[j]
		}
		for j < m-1 {
			j++
			res += a[i] + b[j]
		}
		fmt.Fprintln(out, res)
	}
}
