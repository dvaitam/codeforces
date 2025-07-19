package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l, r int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	intervals := make([]interval, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &intervals[i].l, &intervals[i].r)
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].l != intervals[j].l {
			return intervals[i].l < intervals[j].l
		}
		return intervals[i].r < intervals[j].r
	})

	vis := make([]int, n+2)
	for _, iv := range intervals {
		for pos := iv.l; pos <= iv.r; pos++ {
			vis[pos]++
		}
	}
	sum1 := make([]int, n+2)
	sum2 := make([]int, n+2)
	total := 0
	for i := 1; i <= n; i++ {
		if vis[i] > 0 {
			total++
		}
		sum1[i] = sum1[i-1]
		sum2[i] = sum2[i-1]
		if vis[i] == 1 {
			sum1[i]++
		}
		if vis[i] == 2 {
			sum2[i]++
		}
	}

	ans := 0
	for i := 0; i < q; i++ {
		l1, r1 := intervals[i].l, intervals[i].r
		for j := i + 1; j < q; j++ {
			l2, r2 := intervals[j].l, intervals[j].r
			t := 0
			if l2 <= r1 {
				// overlapping segments
				t += sum2[min(r1, r2)] - sum2[l2-1]
				t += sum1[l2-1] - sum1[l1-1]
				maxr := r1
				if r2 > r1 {
					maxr = r2
				}
				t += sum1[maxr] - sum1[min(r1, r2)]
			} else {
				// disjoint segments
				t += sum1[r1] - sum1[l1-1]
				t += sum1[r2] - sum1[l2-1]
			}
			if total-t > ans {
				ans = total - t
			}
		}
	}
	fmt.Fprintln(out, ans)
}
