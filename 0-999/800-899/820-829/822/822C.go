package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type voucher struct {
	l, r int
	c    int
	len  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}

	segs := make([]voucher, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &segs[i].l, &segs[i].r, &segs[i].c)
		segs[i].len = segs[i].r - segs[i].l + 1
	}

	sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
	byR := make([]voucher, n)
	copy(byR, segs)
	sort.Slice(byR, func(i, j int) bool { return byR[i].r < byR[j].r })

	const INF int64 = 1<<63 - 1
	best := make([]int64, x+1)
	for i := range best {
		best[i] = INF
	}

	ans := INF
	j := 0
	for _, s := range segs {
		for j < n && byR[j].r < s.l {
			l := byR[j].len
			if l <= x && int64(byR[j].c) < best[l] {
				best[l] = int64(byR[j].c)
			}
			j++
		}
		other := x - s.len
		if other >= 0 && best[other] != INF {
			cost := int64(s.c) + best[other]
			if cost < ans {
				ans = cost
			}
		}
	}

	if ans == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
