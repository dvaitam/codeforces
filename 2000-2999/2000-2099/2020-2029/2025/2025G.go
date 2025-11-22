package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segTree struct {
	n       int
	lenPref []int64
	min     []int64
	max     []int64
	sum     []int64
	lazy    []int64
}

func newSegTree(lengths []int64) *segTree {
	n := len(lengths)
	lenPref := make([]int64, n+1)
	for i, v := range lengths {
		lenPref[i+1] = lenPref[i] + v
	}
	size := 4 * n
	return &segTree{
		n:       n,
		lenPref: lenPref,
		min:     make([]int64, size),
		max:     make([]int64, size),
		sum:     make([]int64, size),
		lazy:    make([]int64, size),
	}
}

func (s *segTree) segLen(l, r int) int64 {
	return s.lenPref[r+1] - s.lenPref[l]
}

// sign: +1 if all values in node are non-negative before the update,
// -1 if all are non-positive.
func (s *segTree) apply(idx, l, r int, delta int64, sign int64) {
	length := s.segLen(l, r)
	s.min[idx] += delta
	s.max[idx] += delta
	if sign > 0 {
		s.sum[idx] += delta * length
	} else {
		s.sum[idx] -= delta * length
	}
	s.lazy[idx] += delta
}

func (s *segTree) push(idx, l, r int) {
	if s.lazy[idx] != 0 && l != r {
		m := (l + r) >> 1
		sign := int64(0)
		if s.min[idx] >= 0 {
			sign = 1
		} else if s.max[idx] <= 0 {
			sign = -1
		} else {
			panic("push on mixed sign node")
		}
		d := s.lazy[idx]
		s.apply(idx<<1, l, m, d, sign)
		s.apply(idx<<1|1, m+1, r, d, sign)
		s.lazy[idx] = 0
	}
}

func (s *segTree) pull(idx int) {
	left := idx << 1
	right := left | 1
	s.min[idx] = min64(s.min[left], s.min[right])
	s.max[idx] = max64(s.max[left], s.max[right])
	s.sum[idx] = s.sum[left] + s.sum[right]
}

func (s *segTree) update(L, R int, delta int64) {
	if L > R || s.n == 0 {
		return
	}
	s.updateRec(1, 0, s.n-1, L, R, delta)
}

func (s *segTree) updateRec(idx, l, r, L, R int, delta int64) {
	if R < l || r < L {
		return
	}
	if L <= l && r <= R {
		newMin := s.min[idx] + delta
		newMax := s.max[idx] + delta
		if s.min[idx] >= 0 && newMin >= 0 {
			s.apply(idx, l, r, delta, 1)
			return
		}
		if s.max[idx] <= 0 && newMax <= 0 {
			s.apply(idx, l, r, delta, -1)
			return
		}
	}
	if l == r {
		s.min[idx] += delta
		s.max[idx] += delta
		s.lazy[idx] = 0
		segLen := s.lenPref[l+1] - s.lenPref[l]
		s.sum[idx] = abs64(s.min[idx]) * segLen
		return
	}
	s.push(idx, l, r)
	m := (l + r) >> 1
	s.updateRec(idx<<1, l, m, L, R, delta)
	s.updateRec(idx<<1|1, m+1, r, L, R, delta)
	s.pull(idx)
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}

	type query struct {
		t int
		v int
	}

	queries := make([]query, q)
	values := make([]int, 0, q+1)
	values = append(values, 0)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i].t, &queries[i].v)
		values = append(values, queries[i].v)
	}

	sort.Ints(values)
	values = unique(values)

	lengths := make([]int64, len(values)-1)
	for i := 0; i+1 < len(values); i++ {
		lengths[i] = int64(values[i+1] - values[i])
	}
	seg := newSegTree(lengths)

	var sumH, sumD int64
	for i := 0; i < q; i++ {
		v := queries[i].v
		pos := sort.SearchInts(values, v)
		if pos > 0 {
			delta := int64(1)
			if queries[i].t == 2 {
				delta = -1
			}
			seg.update(0, pos-1, delta)
		}
		if queries[i].t == 1 {
			sumH += int64(v)
		} else {
			sumD += int64(v)
		}
		area := int64(0)
		if seg.n > 0 {
			area = seg.sum[1]
		}
		ans := (3*sumH + sumD - area) / 2
		fmt.Fprintln(out, ans)
	}
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	w := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[w] = a[i]
			w++
		}
	}
	return a[:w]
}
