package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, val int) {
	for idx <= f.n {
		f.tree[idx] ^= val
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int {
	res := 0
	for idx > 0 {
		res ^= f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	parity := make([]int, m+1)
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(in, &c)
		parity[c] ^= 1
	}

	prefCount := make([]int, m+1)
	const maxB = 18
	prefBit := make([][]int, maxB)
	for b := 0; b < maxB; b++ {
		prefBit[b] = make([]int, m+1)
	}
	for i := 1; i <= m; i++ {
		prefCount[i] = prefCount[i-1] ^ parity[i]
		if parity[i] == 1 {
			for b := 0; b < maxB; b++ {
				if (i>>b)&1 == 1 {
					prefBit[b][i] = prefBit[b][i-1] ^ 1
				} else {
					prefBit[b][i] = prefBit[b][i-1]
				}
			}
		} else {
			for b := 0; b < maxB; b++ {
				prefBit[b][i] = prefBit[b][i-1]
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	L := make([]int, q)
	R := make([]int, q)
	countParity := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &L[i], &R[i])
		countParity[i] = prefCount[R[i]] ^ prefCount[L[i]-1]
	}

	ans := make([]int, q)

	for b := 0; b < maxB; b++ {
		mod := 1 << b
		// build positions by residue
		pos := make([][]int, mod)
		for i := 1; i <= m; i++ {
			if parity[i] == 1 {
				r := i & (mod - 1)
				pos[r] = append(pos[r], i)
			}
		}
		// prepare queries for this bit
		type qInfo struct{ t, l, r, idx int }
		qb := make([]qInfo, q)
		for i := 0; i < q; i++ {
			qb[i] = qInfo{t: L[i] & (mod - 1), l: L[i], r: R[i], idx: i}
		}
		sort.Slice(qb, func(i, j int) bool { return qb[i].t < qb[j].t })
		fw := NewFenwick(m)
		cur := 0
		for _, qu := range qb {
			for cur < qu.t {
				for _, p := range pos[cur] {
					fw.Add(p, 1)
				}
				cur++
			}
			borrow := fw.Sum(qu.r) ^ fw.Sum(qu.l-1)
			bitParity := (prefBit[b][qu.r] ^ prefBit[b][qu.l-1]) ^ (((qu.l >> b) & 1) * countParity[qu.idx]) ^ borrow
			if bitParity&1 == 1 {
				ans[qu.idx] |= 1 << b
			}
		}
	}

	for i := 0; i < q; i++ {
		if ans[i] != 0 {
			fmt.Fprint(out, "A")
		} else {
			fmt.Fprint(out, "B")
		}
	}
	out.WriteByte('\n')
}
