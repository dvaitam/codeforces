package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct{ l, r int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var x int
	fmt.Fscan(in, &x)

	intervals := make([]Interval, n)
	for i := 0; i < n; i++ {
		var tl, tr, l, r int
		fmt.Fscan(in, &tl, &tr, &l, &r)
		a := l - tr
		b := r - tl
		if a > b {
			a, b = b, a
		}
		intervals[i] = Interval{a, b}
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].l == intervals[j].l {
			return intervals[i].r < intervals[j].r
		}
		return intervals[i].l < intervals[j].l
	})
	merged := make([]Interval, 0, n)
	for _, iv := range intervals {
		if len(merged) == 0 || iv.l > merged[len(merged)-1].r {
			merged = append(merged, iv)
		} else {
			if iv.r > merged[len(merged)-1].r {
				merged[len(merged)-1].r = iv.r
			}
		}
	}

	ans := 0
	if len(merged) == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	if x < merged[0].l {
		ans = merged[0].l - x
	} else if x > merged[len(merged)-1].r {
		ans = x - merged[len(merged)-1].r
	} else {
		ans = 0
		found := false
		for i, iv := range merged {
			if x < iv.l {
				ans = iv.l - x
				found = true
				break
			} else if x <= iv.r {
				left := x - iv.l
				right := iv.r - x
				if left < right {
					ans = left
				} else {
					ans = right
				}
				found = true
				break
			} else if i+1 < len(merged) && x < merged[i+1].l {
				ans = 0
				found = true
				break
			}
		}
		if !found {
			ans = 0
		}
	}
	fmt.Fprintln(out, ans)
}
