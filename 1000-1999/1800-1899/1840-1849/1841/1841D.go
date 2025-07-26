package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	l int
	r int
}

func minRemovals(segs []Segment) int {
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].r == segs[j].r {
			return segs[i].l < segs[j].l
		}
		return segs[i].r < segs[j].r
	})
	n := len(segs)
	rvals := make([]int, n)
	for i := 0; i < n; i++ {
		rvals[i] = segs[i].r
	}
	dp := make([]int, n+1)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = prefix[i-1]
		li, ri := segs[i-1].l, segs[i-1].r
		for j := i - 1; j > 0; j-- {
			lj, rj := segs[j-1].l, segs[j-1].r
			if rj < li {
				break
			}
			if ri < lj {
				continue
			}
			start := li
			if lj < start {
				start = lj
			}
			p := sort.SearchInts(rvals, start) - 1
			cand := 2
			if p >= 0 {
				cand = prefix[p+1] + 2
			}
			if cand > dp[i] {
				dp[i] = cand
			}
		}
		if dp[i] > prefix[i-1] {
			prefix[i] = dp[i]
		} else {
			prefix[i] = prefix[i-1]
		}
	}
	return n - prefix[n]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		segs := make([]Segment, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &segs[i].l, &segs[i].r)
		}
		ans := minRemovals(segs)
		fmt.Fprintln(writer, ans)
	}
}
