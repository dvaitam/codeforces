package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

		prefix := make([]int64, n+1) // using largest elements of b
		for t2 := 1; t2 <= n; t2++ {
			diff := abs64(a[t2-1] - b[m-t2])
			prefix[t2] = prefix[t2-1] + diff
		}

		suffix := make([]int64, n+1) // using smallest elements of b
		for k := 1; k <= n; k++ {
			diff := abs64(a[n-k] - b[k-1])
			suffix[k] = suffix[k-1] + diff
		}

		var ans int64
		for t2 := 0; t2 <= n; t2++ {
			total := prefix[t2] + suffix[n-t2]
			if total > ans {
				ans = total
			}
		}
		fmt.Fprintln(out, ans)
	}
	out.Flush()
}