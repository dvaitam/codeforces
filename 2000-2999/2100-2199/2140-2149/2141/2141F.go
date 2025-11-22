package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf int64 = 1 << 60

type segTree struct {
	n   int
	tr  []int64
	ans []int64
}

func newSegTree(n int) *segTree {
	tr := make([]int64, 4*n)
	for i := range tr {
		tr[i] = inf
	}
	return &segTree{n: n, tr: tr, ans: make([]int64, n)}
}

func (s *segTree) update(id, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		if val < s.tr[id] {
			s.tr[id] = val
		}
		return
	}
	mid := (l + r) >> 1
	s.update(id<<1, l, mid, ql, qr, val)
	s.update(id<<1|1, mid+1, r, ql, qr, val)
}

func (s *segTree) collect(id, l, r int, carry int64) {
	if s.tr[id] < carry {
		carry = s.tr[id]
	}
	if l == r {
		s.ans[l] = carry
		return
	}
	mid := (l + r) >> 1
	s.collect(id<<1, l, mid, carry)
	s.collect(id<<1|1, mid+1, r, carry)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		// count frequencies
		cntMap := make(map[int]int)
		for _, v := range a {
			cntMap[v]++
		}
		freqs := make([]int, 0, len(cntMap))
		for _, v := range cntMap {
			freqs = append(freqs, v)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(freqs)))

		d := len(freqs)
		preSum := make([]int64, d+1)
		for i, v := range freqs {
			preSum[i+1] = preSum[i] + int64(v)
		}

		maxF := 0
		for _, f := range freqs {
			if f > maxF {
				maxF = f
			}
		}

		freqCnt := make([]int, maxF+2)
		for _, f := range freqs {
			freqCnt[f]++
		}
		cntGT := make([]int, maxF+2)
		sumGT := make([]int64, maxF+2)
		for i := maxF; i >= 0; i-- {
			cntGT[i] = cntGT[i+1] + freqCnt[i+1]
			sumGT[i] = sumGT[i+1] + int64(i+1)*int64(freqCnt[i+1])
		}

		seg := newSegTree(n)

		for t := 0; t <= maxF; t++ {
			p := cntGT[t] // number of frequencies greater than t
			L := sumGT[t] - int64(p)*int64(t)

			// c = 0 interval [L, n-1]
			if L < int64(n) {
				l := int64(L)
				if l < 0 {
					l = 0
				}
				seg.update(1, 0, n-1, int(l), n-1, int64(t))
			}

			Sprev := int64(0)
			for c := 1; c <= p; c++ {
				Sc := preSum[c] - int64(c)*int64(t)
				l := int64(L - Sc)
				r := int64(L - Sprev - 1)
				if l < 0 {
					l = 0
				}
				if r >= int64(n) {
					r = int64(n - 1)
				}
				if l <= r {
					seg.update(1, 0, n-1, int(l), int(r), int64(t+c))
				}
				Sprev = Sc
			}
		}

		seg.collect(1, 0, n-1, inf)
		for i, v := range seg.ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
