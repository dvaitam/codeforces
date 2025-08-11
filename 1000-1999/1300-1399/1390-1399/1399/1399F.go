package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Seg struct {
	l, r int
}

type Child struct {
	l, r, w int
}

func solve(segs []Seg) int {
	sort.Slice(segs, func(i, j int) bool {
		li := segs[i].r - segs[i].l
		lj := segs[j].r - segs[j].l
		if li != lj {
			return li < lj
		}
		return segs[i].l < segs[j].l
	})
	n := len(segs)
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		var childs []Child
		for j := 0; j < i; j++ {
			if segs[i].l <= segs[j].l && segs[j].r <= segs[i].r {
				childs = append(childs, Child{segs[j].l, segs[j].r, 1 + dp[j]})
			}
		}
		if len(childs) == 0 {
			dp[i] = 0
			continue
		}
		sort.Slice(childs, func(a, b int) bool { return childs[a].r < childs[b].r })
		ends := make([]int, len(childs))
		for k := range childs {
			ends[k] = childs[k].r
		}
		local := make([]int, len(childs)+1)
		for k := 1; k <= len(childs); k++ {
			lCur := childs[k-1].l
			wCur := childs[k-1].w
			q := sort.SearchInts(ends, lCur) - 1
			prev := 0
			if q >= 0 {
				prev = local[q+1]
			}
			take := wCur + prev
			if take > local[k-1] {
				local[k] = take
			} else {
				local[k] = local[k-1]
			}
		}
		dp[i] = local[len(childs)]
	}
	childs := make([]Child, n)
	for k := 0; k < n; k++ {
		childs[k] = Child{segs[k].l, segs[k].r, 1 + dp[k]}
	}
	sort.Slice(childs, func(a, b int) bool { return childs[a].r < childs[b].r })
	ends := make([]int, n)
	for k := range childs {
		ends[k] = childs[k].r
	}
	local := make([]int, n+1)
	for k := 1; k <= n; k++ {
		lCur := childs[k-1].l
		wCur := childs[k-1].w
		q := sort.SearchInts(ends, lCur) - 1
		prev := 0
		if q >= 0 {
			prev = local[q+1]
		}
		take := wCur + prev
		if take > local[k-1] {
			local[k] = take
		} else {
			local[k] = local[k-1]
		}
	}
	return local[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		segs := make([]Seg, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
		}
		fmt.Fprintln(out, solve(segs))
	}
}
