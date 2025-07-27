package main

import (
	"bufio"
	"fmt"
	"os"
)

// This implementation uses a basic segment tree to approximate the solution for
// problemE.txt. The approach follows a greedy strategy where values from the
// permutation are added from largest to smallest until every prefix contains
// strictly more values than bombs. While this implementation does not strictly
// reproduce the official algorithm, it demonstrates one possible way to model
// the problem using range updates on prefixes.

type segTree struct {
	n    int
	tree []int
	lazy []int
}

func newSegTree(size int) *segTree {
	n := 1
	for n < size {
		n <<= 1
	}
	return &segTree{n: n, tree: make([]int, 2*n), lazy: make([]int, 2*n)}
}

func (s *segTree) apply(idx, val int) {
	s.tree[idx] += val
	s.lazy[idx] += val
}

func (s *segTree) push(idx int) {
	if s.lazy[idx] != 0 {
		s.apply(idx*2, s.lazy[idx])
		s.apply(idx*2+1, s.lazy[idx])
		s.lazy[idx] = 0
	}
}

func (s *segTree) addRange(l, r, val, idx, lo, hi int) {
	if l <= lo && hi <= r {
		s.apply(idx, val)
		return
	}
	s.push(idx)
	mid := (lo + hi) / 2
	if l <= mid {
		s.addRange(l, r, val, idx*2, lo, mid)
	}
	if r > mid {
		s.addRange(l, r, val, idx*2+1, mid+1, hi)
	}
	if s.tree[idx*2] < s.tree[idx*2+1] {
		s.tree[idx] = s.tree[idx*2]
	} else {
		s.tree[idx] = s.tree[idx*2+1]
	}
}

func (s *segTree) AddPrefix(r, val int) {
	if r < 1 {
		return
	}
	if r > s.n {
		r = s.n
	}
	s.addRange(1, r, val, 1, 1, s.n)
}

func (s *segTree) Min() int { return s.tree[1] }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
		pos[p[i]] = i
	}
	q := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &q[i])
	}

	seg := newSegTree(n)
	ans := n
	seg.AddPrefix(pos[ans], 1)
	res := make([]int, n)

	for i := 0; i < n; i++ {
		for seg.Min() < 0 && ans > 1 {
			ans--
			seg.AddPrefix(pos[ans], 1)
		}
		res[i] = ans
		seg.AddPrefix(q[i], -1)
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}
