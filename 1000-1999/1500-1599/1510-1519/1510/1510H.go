package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	L, R int
	idx  int
	len  int
}

type Interval struct {
	L, R int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	segs := make([]Segment, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &segs[i].L, &segs[i].R)
		segs[i].idx = i
		segs[i].len = segs[i].R - segs[i].L
	}

	sorted := make([]Segment, n)
	copy(sorted, segs)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].len == sorted[j].len {
			return sorted[i].L < sorted[j].L
		}
		return sorted[i].len < sorted[j].len
	})

	assigned := make([]Interval, 0, n)
	ans := make([]Interval, n)

	for _, seg := range sorted {
		L, R := seg.L, seg.R
		inside := make([]Interval, 0)
		for _, iv := range assigned {
			if iv.L >= L && iv.R <= R {
				inside = append(inside, iv)
			}
		}
		sort.Slice(inside, func(i, j int) bool { return inside[i].L < inside[j].L })
		bestLen := -1
		bestStart := L
		prev := L
		for _, iv := range inside {
			if iv.L > prev {
				curLen := iv.L - prev
				if curLen > bestLen {
					bestLen = curLen
					bestStart = prev
				}
			}
			if iv.R > prev {
				prev = iv.R
			}
		}
		if R > prev {
			curLen := R - prev
			if curLen > bestLen {
				bestLen = curLen
				bestStart = prev
			}
		}
		if bestLen <= 0 {
			bestLen = 1
			bestStart = L
		}
		interval := Interval{bestStart, bestStart + bestLen}
		ans[seg.idx] = interval
		assigned = append(assigned, interval)
	}

	total := 0
	for _, iv := range ans {
		total += iv.R - iv.L
	}
	fmt.Fprintln(out, total)
	for _, iv := range ans {
		fmt.Fprintf(out, "%d %d\n", iv.L, iv.R)
	}
}
