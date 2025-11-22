package main

import (
	"bufio"
	"fmt"
	"os"
)

type Solver struct {
	n    int
	l, r []int64
	p, q []int64
	vis  []int
	ver  int
}

func (s *Solver) check(dy int, orient int) bool {
	// orient: 0 -> A is on the left (touches left boundary), 1 -> A on the right.
	n := s.n
	if dy == 0 {
		// Need all widths equal and even; additionally A must be connected.
		w := s.r[0] - s.l[0] + 1
		if w%2 != 0 {
			return false
		}
		for i := 1; i < n; i++ {
			if s.r[i]-s.l[i]+1 != w {
				return false
			}
		}
		dx := w / 2
		// Try both orientations: A on the left half or on the right half.
		for orient := 0; orient < 2; orient++ {
			ok := true
			prevL := s.l[0]
			prevR := prevL + dx - 1
			if orient == 1 {
				prevL = s.l[0] + dx
				prevR = s.r[0]
			}
			for i := 1; i < n && ok; i++ {
				curL := s.l[i]
				curR := curL + dx - 1
				if orient == 1 {
					curL = s.l[i] + dx
					curR = s.r[i]
				}
				// Consecutive rows must overlap to keep connectivity.
				if curL > prevR || prevL > curR {
					ok = false
					break
				}
				prevL, prevR = curL, curR
			}
			if ok {
				return true
			}
		}
		return false
	}

	if dy > n/2 {
		return false
	}

	m := n - dy
	s.ver++
	ver := s.ver
	// Initialize rows that only belong to A (top dy rows).
	limit := dy
	if limit > m {
		limit = m
	}
	for i := 0; i < limit; i++ {
		s.p[i] = s.l[i]
		s.q[i] = s.r[i]
		s.vis[i] = ver
	}

	dxSet := false
	var dx int64

	// Overlap rows: both A(row i) and B(row i-dy) present.
	for j := dy; j < n-dy; j++ {
		i := j
		prev := i - dy
		if s.vis[prev] != ver {
			return false
		}
		pp := s.p[prev]
		qq := s.q[prev]

		if orient == 0 {
			need := s.r[j] - qq
			if !dxSet {
				dx, dxSet = need, true
			} else if dx != need {
				return false
			}
			s.p[i] = s.l[j]
			s.q[i] = pp + dx - 1
			if qq+dx != s.r[j] {
				return false
			}
		} else {
			need := s.l[j] - pp
			if !dxSet {
				dx, dxSet = need, true
			} else if dx != need {
				return false
			}
			s.p[i] = qq + dx + 1
			s.q[i] = s.r[j]
			if pp+dx != s.l[j] {
				return false
			}
		}
		s.vis[i] = ver
	}

	// Rows that only belong to B (bottom dy rows): j from n-dy to n-1
	for j := n - dy; j < n; j++ {
		i := j - dy
		if !dxSet {
			if s.vis[i] == ver {
				dx = s.l[j] - s.p[i]
			} else {
				dx = s.l[j] - s.l[i]
			}
			dxSet = true
		}

		if s.vis[i] != ver {
			s.p[i] = s.l[j] - dx
			s.q[i] = s.r[j] - dx
			s.vis[i] = ver
		} else {
			pi := s.l[j] - dx
			qi := s.r[j] - dx
			if s.p[i] != pi || s.q[i] != qi {
				return false
			}
		}
		// Also ensure right boundary aligns.
		if s.r[j]-s.q[i] != dx {
			return false
		}
	}

	return dxSet
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		l := make([]int64, n)
		r := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
		}

		s := Solver{
			n:   n,
			l:   l,
			r:   r,
			p:   make([]int64, n),
			q:   make([]int64, n),
			vis: make([]int, n),
		}

		ok := s.check(0, 0)
		if !ok {
			for dy := 1; dy < n && !ok; dy++ {
				if s.check(dy, 0) || s.check(dy, 1) {
					ok = true
				}
			}
		}

		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
