package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Simple segment tree for range minimum queries
// Works with int64 values.
const inf int64 = 1<<63 - 1

type segTree struct {
	n    int
	tree []int64
}

func newSegTree(arr []int64) *segTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	tree := make([]int64, 2*n)
	for i := range tree {
		tree[i] = inf
	}
	for i, v := range arr {
		tree[n+i] = v
	}
	for i := n - 1; i > 0; i-- {
		if tree[i<<1] < tree[i<<1|1] {
			tree[i] = tree[i<<1]
		} else {
			tree[i] = tree[i<<1|1]
		}
	}
	return &segTree{n: n, tree: tree}
}

func (s *segTree) query(l, r int) int64 {
	if l > r {
		return inf
	}
	l += s.n
	r += s.n
	res := inf
	for l <= r {
		if l&1 == 1 {
			if s.tree[l] < res {
				res = s.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if s.tree[r] < res {
				res = s.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var a, b, c, d, start, length int64
	if _, err := fmt.Fscan(in, &n, &a, &b, &c, &d, &start, &length); err != nil {
		return
	}

	times := make([]int64, n)
	typ := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &times[i], &typ[i])
	}

	before := make([]int64, n)
	after := make([]int64, n)
	for i := 0; i < n; i++ {
		if typ[i] == 1 {
			before[i] = a
			after[i] = c
		} else {
			before[i] = -b
			after[i] = -d
		}
	}

	prefBefore := make([]int64, n+1)
	prefAfter := make([]int64, n+1)
	minBefore := make([]int64, n+1)
	prefBefore[0] = start
	prefAfter[0] = start
	minBefore[0] = start
	for i := 0; i < n; i++ {
		prefBefore[i+1] = prefBefore[i] + before[i]
		prefAfter[i+1] = prefAfter[i] + after[i]
		if prefBefore[i+1] < minBefore[i] {
			minBefore[i+1] = prefBefore[i+1]
		} else {
			minBefore[i+1] = minBefore[i]
		}
	}

	diffPrefix := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		diffPrefix[i] = prefAfter[i] - prefBefore[i]
	}

	seg := newSegTree(prefAfter)

	candMap := map[int64]struct{}{0: {}}
	for i := 0; i < n; i++ {
		t1 := times[i] - length + 1
		if t1 >= 0 {
			candMap[t1] = struct{}{}
		}
		candMap[times[i]+1] = struct{}{}
	}
	cands := make([]int64, 0, len(candMap))
	for v := range candMap {
		if v >= 0 {
			cands = append(cands, v)
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i] < cands[j] })

	for _, t := range cands {
		L := sort.Search(len(times), func(i int) bool { return times[i] >= t })
		Rp := sort.Search(len(times), func(i int) bool { return times[i] >= t+length })
		R := Rp - 1
		if minBefore[L] < 0 {
			continue
		}
		if L <= R {
			segMin := seg.query(L+1, R+1) - diffPrefix[L]
			if segMin < 0 {
				continue
			}
		}
		fmt.Fprintln(out, t)
		return
	}
	fmt.Fprintln(out, -1)
}
