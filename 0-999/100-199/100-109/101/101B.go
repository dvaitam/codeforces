package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 1000000007

type seg struct {
	r, l int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	segs := make([]seg, m)
	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		segs[i] = seg{r: r, l: l}
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].r < segs[j].r })
	// build endpoints with sentinel 0
	b := []int{0}
	c := make([]int, m)
	last := -1
	for i, s := range segs {
		if s.r != last {
			last = s.r
			b = append(b, s.r)
		}
		c[i] = len(b) - 1
	}
	E := len(b) - 1
	if b[E] != n {
		fmt.Println(0)
		return
	}
	f := make([]int, E+1)
	ssum := make([]int, E+1)
	f[0], ssum[0] = 1, 1
	for i, seg := range segs {
		ci := c[i]
		// find lowest index lo where b[lo] >= seg.l
		lo, hi := 0, ci
		for lo < hi {
			mid := (lo + hi) / 2
			if b[mid] < seg.l {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo < ci {
			add := ssum[ci-1]
			if lo > 0 {
				add = (add - ssum[lo-1] + mod) % mod
			}
			f[ci] = (f[ci] + add) % mod
		}
		// update prefix sum when next segment has different endpoint or last
		if i+1 == m || c[i+1] != ci {
			ssum[ci] = (ssum[ci-1] + f[ci]) % mod
		}
	}
	fmt.Println(f[E])
}
