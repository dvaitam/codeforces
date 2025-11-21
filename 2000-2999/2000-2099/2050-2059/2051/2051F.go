package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l int64
	r int64
}

func applyIntervals(intervals []interval, n, a int64) []interval {
	newIntervals := make([]interval, 0, len(intervals)*3)
	for _, seg := range intervals {
		L, R := seg.l, seg.r

		if L <= a-1 {
			r1 := R
			if r1 > a-1 {
				r1 = a - 1
			}
			if L <= r1 {
				newIntervals = append(newIntervals, interval{L, r1 + 1})
			}
		}

		if R >= a+1 {
			l2 := L
			if l2 < a+1 {
				l2 = a + 1
			}
			if l2 <= R {
				newIntervals = append(newIntervals, interval{l2 - 1, R})
			}
		}

		if L <= a && a <= R {
			newIntervals = append(newIntervals, interval{1, 1})
			newIntervals = append(newIntervals, interval{n, n})
		}
	}

	if len(newIntervals) == 0 {
		return newIntervals
	}

	sort.Slice(newIntervals, func(i, j int) bool {
		if newIntervals[i].l == newIntervals[j].l {
			return newIntervals[i].r < newIntervals[j].r
		}
		return newIntervals[i].l < newIntervals[j].l
	})

	merged := make([]interval, 0, len(newIntervals))
	for _, cur := range newIntervals {
		if len(merged) == 0 || cur.l > merged[len(merged)-1].r+1 {
			merged = append(merged, cur)
		} else {
			if cur.r > merged[len(merged)-1].r {
				merged[len(merged)-1].r = cur.r
			}
		}
	}

	return merged
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n, m int64
		var q int
		fmt.Fscan(in, &n, &m, &q)

		intervals := []interval{{m, m}}

		for i := 0; i < q; i++ {
			var a int64
			fmt.Fscan(in, &a)
			intervals = applyIntervals(intervals, n, a)

			var total int64
			for _, seg := range intervals {
				total += seg.r - seg.l + 1
			}

			if i+1 == q {
				fmt.Fprintln(out, total)
			} else {
				fmt.Fprint(out, total, " ")
			}
		}
	}
}

