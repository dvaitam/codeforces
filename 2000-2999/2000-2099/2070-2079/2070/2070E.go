package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Fenwick tree for prefix sums.
type fenwick struct {
	n int
	f []int64
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, f: make([]int64, n+1)}
}

func (fw *fenwick) add(idx int, delta int64) {
	for idx <= fw.n {
		fw.f[idx] += delta
		idx += idx & -idx
	}
}

func (fw *fenwick) sum(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += fw.f[idx]
		idx -= idx & -idx
	}
	return res
}

func upperBound(a []int, x int) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] <= x {
			l = m + 1
		} else {
			r = m
		}
	}
	return l
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	pref1 := make([]int, n+1)
	for i := 0; i < n; i++ {
		pref1[i+1] = pref1[i]
		if s[i] == '1' {
			pref1[i+1]++
		}
	}

	// A[i] = i - 4 * pref1[i]
	A := make([]int, n+1)
	for i := 0; i <= n; i++ {
		A[i] = i - 4*pref1[i]
	}

	// Prepare compression for Fenwick queries (difference >= 2 case).
	coords := make([]int, n+1)
	copy(coords, A)
	sort.Ints(coords)
	coords = unique(coords)

	fw := newFenwick(len(coords))

	total := int64(0)
	// Counting substrings where A[r]-A[l] >= 2 (i.e., len - 4*ones >= 2).
	for _, v := range A {
		thr := v - 2
		pos := upperBound(coords, thr)
		if pos > 0 {
			total += fw.sum(pos)
		}
		// update with current prefix
		idx := upperBound(coords, v) // first greater than v, so index of v is idx
		fw.add(idx, 1)
	}

	// Counting substrings with A[r]-A[l] == -1 (special winning case).
	freq := make(map[int]int64)
	for _, v := range A {
		total += freq[v+1]
		freq[v]++
	}

	fmt.Fprintln(out, total)
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	res := a[:1]
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			res = append(res, a[i])
		}
	}
	return res
}
