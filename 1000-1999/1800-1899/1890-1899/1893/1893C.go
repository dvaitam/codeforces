package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Pair represents a set index and the count of a particular value in that set.
type Pair struct {
	idx int
	cnt int64
}

// Set holds parameters for each multiset.
type Set struct {
	l, r, total int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		sets := make([]Set, m)
		valMap := make(map[int64][]Pair)
		uniq := make(map[int64]struct{})
		var L, R int64
		for i := 0; i < m; i++ {
			var n int
			var l, r int64
			fmt.Fscan(reader, &n, &l, &r)
			L += l
			R += r
			a := make([]int64, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &a[j])
			}
			c := make([]int64, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &c[j])
			}
			var total int64
			for j := 0; j < n; j++ {
				total += c[j]
				valMap[a[j]] = append(valMap[a[j]], Pair{idx: i, cnt: c[j]})
				uniq[a[j]] = struct{}{}
			}
			sets[i] = Set{l: l, r: r, total: total}
		}

		values := make([]int64, 0, len(uniq))
		for v := range uniq {
			values = append(values, v)
		}
		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

		prev := L - 1
		gap := false
		for _, v := range values {
			if v < L {
				continue
			}
			if v > R {
				break
			}
			if v > prev+1 {
				gap = true
				break
			}
			prev = v
		}
		if !gap && prev < R {
			gap = true
		}
		if gap {
			fmt.Fprintln(writer, 0)
			continue
		}

		baseCap := R - L
		best := int64(1<<63 - 1)
		for _, S := range values {
			if S < L || S > R {
				continue
			}
			F0 := int64(0)
			nonSCap := baseCap
			extraCap := int64(0)
			for _, p := range valMap[S] {
				set := sets[p.idx]
				t := set.total - p.cnt
				base := int64(0)
				if t < set.l {
					base = set.l - t
				}
				minrt := set.r
				if t < minrt {
					minrt = t
				}
				capi := int64(0)
				if minrt > set.l {
					capi = minrt - set.l
				}
				nonSCap += capi - (set.r - set.l)
				extra := set.r - set.l - capi
				remain := p.cnt - base
				if remain < extra {
					extra = remain
				}
				if extra < 0 {
					extra = 0
				}
				extraCap += extra
				F0 += base
			}
			need := int64(0)
			if S-L > nonSCap {
				need = S - L - nonSCap
			}
			if need <= extraCap {
				cand := F0 + need
				if cand < best {
					best = cand
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
