package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Suffix array construction for integer alphabet
func buildSuffixArray(s []int) []int {
	n := len(s)
	sa := make([]int, n)
	for i := 0; i < n; i++ {
		sa[i] = i
	}
	// rank and tmp
	rank := make([]int, n)
	tmp := make([]int, n)
	for i := 0; i < n; i++ {
		rank[i] = s[i]
	}
	for k := 1; k < n; k <<= 1 {
		// compare by (rank[i], rank[i+k])
		cmp := func(i, j int) bool {
			if rank[i] != rank[j] {
				return rank[i] < rank[j]
			}
			si := -1
			if i+k < n {
				si = rank[i+k]
			}
			sj := -1
			if j+k < n {
				sj = rank[j+k]
			}
			return si < sj
		}
		sort.Slice(sa, func(i, j int) bool {
			return cmp(sa[i], sa[j])
		})
		tmp[sa[0]] = 0
		for i := 1; i < n; i++ {
			tmp[sa[i]] = tmp[sa[i-1]]
			if cmp(sa[i-1], sa[i]) {
				tmp[sa[i]]++
			}
		}
		copy(rank, tmp)
		if rank[sa[n-1]] == n-1 {
			break
		}
	}
	return sa
}

// LCP array Kasai's algorithm
func buildLCP(s []int, sa []int) []int {
	n := len(s)
	rank := make([]int, n)
	for i, p := range sa {
		rank[p] = i
	}
	lcp := make([]int, n)
	h := 0
	for i := 0; i < n; i++ {
		if rank[i] > 0 {
			j := sa[rank[i]-1]
			for i+h < n && j+h < n && s[i+h] == s[j+h] {
				h++
			}
			lcp[rank[i]] = h
			if h > 0 {
				h--
			}
		}
	}
	return lcp
}

// BIT for ints
type BIT struct {
	n int
	bit []int
}

func NewBIT(n int) *BIT {
	return &BIT{n, make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
	i++
	for ; i <= b.n; i += i & -i {
		b.bit[i] += v
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	i++
	if i > b.n {
		i = b.n
	}
	for ; i > 0; i -= i & -i {
		s += b.bit[i]
	}
	return s
}

func (b *BIT) RangeSum(l, r int) int {
	if l > r {
		return 0
	}
	return b.Sum(r) - b.Sum(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   var str string
   fmt.Fscan(in, &str)
   sraw := []byte(str)
	// build P and left_less
	P := make([]int, n+1)
	for i := 1; i <= n; i++ {
		P[i] = P[i-1]
		if sraw[i-1] == '(' {
			P[i]++
		} else {
			P[i]--
		}
	}
	// left_less: previous index with P[j] < P[i]
	left_less := make([]int, n+1)
	s := make([]int, n+1)
	for i := 0; i <= n; i++ {
		s[i] = P[i]
	}
	stack := make([]int, 0, n+1)
	for i := 0; i <= n; i++ {
		for len(stack) > 0 && s[stack[len(stack)-1]] >= s[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			left_less[i] = -1
		} else {
			left_less[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	// build suffix array on sraw
	// map '('->1, ')'->2
	s2 := make([]int, n)
	for i := 0; i < n; i++ {
		if sraw[i] == '(' {
			s2[i] = 1
		} else {
			s2[i] = 2
		}
	}
	sa := buildSuffixArray(s2)
	lcp := buildLCP(s2, sa)
	// compress P values
	vals := make([]int, n+1)
	copy(vals, P)
	sort.Ints(vals)
	uvals := vals[:1]
	for i := 1; i < len(vals); i++ {
		if vals[i] != vals[i-1] {
			uvals = append(uvals, vals[i])
		}
	}
	m := len(uvals)
	comp := make(map[int]int, m)
	for i, v := range uvals {
		comp[v] = i
	}
	// Ppos, events, queries per group
	Ppos := make([][]int, m)
	events := make([][]struct{time, idx int}, m)
	queries := make([][]struct{time, r0 int}, m)
	for i := 0; i <= n; i++ {
		ci := comp[P[i]]
		Ppos[ci] = append(Ppos[ci], i)
	}
	// For each r from 1..n, create event at time=max(left_less[r]+1,1)
	for ci := 0; ci < m; ci++ {
		events[ci] = []struct{time, idx int}{}
		queries[ci] = []struct{time, r0 int}{}
	}
	// map pos to index in Ppos per group
	posIndex := make([]int, n+1)
	for ci := 0; ci < m; ci++ {
		for j, pos := range Ppos[ci] {
			posIndex[pos] = j
		}
	}
	for r := 1; r <= n; r++ {
		t := left_less[r] + 1
		if t < 1 {
			t = 1
		}
		ci := comp[P[r]]
		idx := posIndex[r]
		events[ci] = append(events[ci], struct{time, idx int}{t, idx})
	}
	// build queries from SA and LCP
	for k, si := range sa {
		pos := si + 1
		h := lcp[k]
		r0 := pos + h
		if r0 > n {
			continue
		}
		ci := comp[P[pos-1]]
		queries[ci] = append(queries[ci], struct{time, r0 int}{pos, r0})
	}
	// answer
	var ans int64
	// process per group
	for ci := 0; ci < m; ci++ {
		sz := len(Ppos[ci])
		if sz == 0 {
			continue
		}
		// sort events by time, queries by time
		e := events[ci]
		sr := queries[ci]
		sort.Slice(e, func(i, j int) bool { return e[i].time < e[j].time })
		sort.Slice(sr, func(i, j int) bool { return sr[i].time < sr[j].time })
		bit := NewBIT(sz)
		ep, qp := 0, 0
		nq, ne := len(sr), len(e)
		for qp < nq {
			t := sr[qp].time
			for ep < ne && e[ep].time <= t {
				bit.Add(e[ep].idx, 1)
				ep++
			}
			// process all queries at this time
			for qp < nq && sr[qp].time == t {
				r0 := sr[qp].r0
				// find first index j0 in Ppos >= r0
				j0 := sort.Search(sz, func(i int) bool { return Ppos[ci][i] >= r0 })
				if j0 < sz {
					cnt := bit.RangeSum(j0, sz-1)
					ans += int64(cnt)
				}
				qp++
			}
			// Next query may have larger time; if no more queries at this t, continue loop
		}
	}
	fmt.Fprintln(out, ans)
}
