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
		m := 2 * n
		times := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &times[i])
		}
		segments := 2*n - 1
		prefix := make([]int64, segments+1)
		prefixWeighted := make([]int64, segments+1)
		prefixOdd := make([]int64, segments+1)
		prefixEven := make([]int64, segments+1)
		prefixTail := make([]int64, segments+1)

		for i := 1; i <= segments; i++ {
			diff := times[i] - times[i-1]
			prefix[i] = prefix[i-1] + diff
			prefixWeighted[i] = prefixWeighted[i-1] + int64(i)*diff
			prefixTail[i] = prefixTail[i-1] + int64(2*n-i)*diff
			if i%2 == 1 {
				prefixOdd[i] = prefixOdd[i-1] + diff
				prefixEven[i] = prefixEven[i-1]
			} else {
				prefixEven[i] = prefixEven[i-1] + diff
				prefixOdd[i] = prefixOdd[i-1]
			}
		}

		sumRange := func(pref []int64, l, r int) int64 {
			if l > r {
				return 0
			}
			return pref[r] - pref[l-1]
		}
		sumParity := func(l, r int, parity int) int64 {
			if l > r {
				return 0
			}
			if parity%2 == 1 {
				return prefixOdd[r] - prefixOdd[l-1]
			}
			return prefixEven[r] - prefixEven[l-1]
		}

		totalTail := prefixTail[segments]
		ans := make([]int64, n)
		for k := 1; k <= n; k++ {
			ramp := prefixWeighted[k]
			var middle int64
			if k < n {
				left := k + 1
				right := 2*n - k
				if left <= right {
					totalMid := sumRange(prefix, left, right)
					sameParity := sumParity(left, right, k)
					diffParity := totalMid - sameParity
					middle = int64(k)*sameParity + int64(k-1)*diffParity
				}
			}
			var tail int64
			if k > 1 {
				start := 2*n - k
				tail = totalTail - prefixTail[start]
			}
			ans[k-1] = ramp + middle + tail
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
