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
	fmt.Fscan(in, &t)
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
		sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
		cur := int64(-1 << 60)
		ok := true
		for i := 0; i < n; i++ {
			best := int64(1 << 62)
			if a[i] >= cur {
				best = a[i]
			}
			target := cur + a[i]
			idx := lowerBound(b, target)
			if idx < len(b) {
				candidate := b[idx] - a[i]
				if candidate >= cur && candidate < best {
					best = candidate
				}
			}
			if best == int64(1<<62) {
				ok = false
				break
			}
			cur = best
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func lowerBound(arr []int64, target int64) int {
	l, r := 0, len(arr)
	for l < r {
		mid := (l + r) >> 1
		if arr[mid] < target {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l
}
