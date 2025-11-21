package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func countPairsGreater(a []int, limit int, threshold int) int64 {
	if limit < 2 {
		return 0
	}
	l, r := 0, limit-1
	var res int64
	for l < r {
		if a[l]+a[r] > threshold {
			res += int64(r - l)
			r--
		} else {
			l++
		}
	}
	return res
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		var ans int64

		// Case 1: triples not containing the largest element.
		if n >= 4 {
			for k := 2; k <= n-2; k++ {
				T := 2 * a[k]
				if a[n-1] > T {
					T = a[n-1]
				}
				threshold := T - a[k]
				ans += countPairsGreater(a, k, threshold)
			}
		}

		// Case 2: triples containing a[n-1] but not a[n-2].
		if n >= 4 {
			ans += countPairsGreater(a, n-2, a[n-1])
		}

		// Case 3: triples containing a[n-1] and a[n-2], but not a[n-3].
		if n >= 4 {
			limit := n - 3 // use elements a[0..n-4]
			if limit > 0 {
				threshold := a[n-1] - a[n-2]
				idx := sort.Search(limit, func(i int) bool { return a[i] > threshold })
				ans += int64(limit - idx)
			}
		}

		// Case 4: triple consisting of the three largest elements.
		if n >= 3 && a[n-3]+a[n-2] > a[n-1] {
			ans++
		}

		fmt.Fprintln(out, ans)
	}
}

