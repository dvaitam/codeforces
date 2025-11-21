package main

import (
	"bufio"
	"fmt"
	"os"
)

type Str struct {
	dbl  []byte
	pref []uint64
}

var (
	n, m int
	base uint64 = 1469598103934665603
	pow  []uint64
	arr  []Str
)

func buildPow(limit int) {
	pow = make([]uint64, limit+1)
	pow[0] = 1
	for i := 1; i <= limit; i++ {
		pow[i] = pow[i-1] * base
	}
}

func buildPref(s string) Str {
	dbl := make([]byte, 2*m)
	copy(dbl, s)
	copy(dbl[m:], s)
	pref := make([]uint64, 2*m+1)
	for i := 0; i < 2*m; i++ {
		val := uint64(dbl[i]-'0') + 1
		pref[i+1] = pref[i]*base + val
	}
	return Str{dbl: dbl, pref: pref}
}

func hash(st *Str, l, r int) uint64 {
	return st.pref[r] - st.pref[l]*pow[r-l]
}

func cmp(i, j, start int) int {
	si := &arr[i]
	sj := &arr[j]
	l, r := 0, m
	for l < r {
		mid := (l + r + 1) >> 1
		if hash(si, start, start+mid) == hash(sj, start, start+mid) {
			l = mid
		} else {
			r = mid - 1
		}
	}
	if l == m {
		return 0
	}
	ci := si.dbl[start+l]
	cj := sj.dbl[start+l]
	if ci > cj {
		return 1
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m)
	buildPow(2 * m)
	arr = make([]Str, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		arr[i] = buildPref(s)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for start := 0; start < m; start++ {
		best := 0
		for i := 1; i < n; i++ {
			if cmp(i, best, start) > 0 {
				best = i
			}
		}
		if start > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, best+1)
	}
	fmt.Fprintln(out)
}
