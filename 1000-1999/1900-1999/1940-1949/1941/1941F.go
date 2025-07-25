package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// existsSumInRange checks if there exist d in D and f in F such that L <= d+f <= R.
func existsSumInRange(D, F []int64, L, R int64) bool {
	i, j := 0, len(F)-1
	for i < len(D) && j >= 0 {
		s := D[i] + F[j]
		if s < L {
			i++
		} else if s > R {
			j--
		} else {
			return true
		}
	}
	return false
}

// can checks whether it's possible to have maximum gap <= limit
// after inserting at most one element with value from D+F.
func can(a, D, F []int64, limit int64) bool {
	pos := -1
	for i := 1; i < len(a); i++ {
		gap := a[i] - a[i-1]
		if gap > limit {
			if pos != -1 {
				return false
			}
			pos = i
		}
	}
	if pos == -1 {
		return true
	}
	gap := a[pos] - a[pos-1]
	if gap > 2*limit {
		return false
	}
	left := maxInt64(a[pos-1], a[pos]-limit)
	right := minInt64(a[pos], a[pos-1]+limit)
	return existsSumInRange(D, F, left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)

		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		d := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &d[i])
		}
		f := make([]int64, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &f[i])
		}

		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
		sort.Slice(f, func(i, j int) bool { return f[i] < f[j] })

		low, high := int64(0), int64(4000000000)
		for low < high {
			mid := (low + high) / 2
			if can(a, d, f, mid) {
				high = mid
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintln(out, low)
	}
}
