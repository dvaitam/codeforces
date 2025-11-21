package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type event struct {
	pos   int64
	delta int
}

func feasible(n int, m int64, k int, h []int64, x []int64, attacks int64) bool {
	evs := make([]event, 0, 2*n)
	for i := 0; i < n; i++ {
		req := (h[i] + attacks - 1) / attacks
		if req > m {
			continue
		}
		radius := m - req
		l := x[i] - radius
		r := x[i] + radius
		evs = append(evs, event{pos: l, delta: 1})
		evs = append(evs, event{pos: r + 1, delta: -1})
	}
	if len(evs) == 0 {
		return false
	}
	sort.Slice(evs, func(i, j int) bool {
		if evs[i].pos == evs[j].pos {
			return evs[i].delta < evs[j].delta
		}
		return evs[i].pos < evs[j].pos
	})
	cur := 0
	for _, e := range evs {
		cur += e.delta
		if cur >= k {
			return true
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		var m int64
		fmt.Fscan(in, &n, &m, &k)
		h := make([]int64, n)
		var maxH int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
			if h[i] > maxH {
				maxH = h[i]
			}
		}
		x := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}

		if !feasible(n, m, k, h, x, maxH) {
			fmt.Fprintln(out, -1)
			continue
		}

		lo, hi := int64(1), maxH
		ans := maxH
		for lo <= hi {
			mid := (lo + hi) / 2
			if feasible(n, m, k, h, x, mid) {
				ans = mid
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(out, ans)
	}
}
