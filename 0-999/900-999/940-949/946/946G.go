package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Fenwick tree storing pairs (length,count) for LIS computation.
type Fenwick struct {
	len []int
	cnt []float64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+2), make([]float64, n+2)}
}

func (f *Fenwick) update(i, l int, c float64) {
	for i < len(f.len) {
		if l > f.len[i] {
			f.len[i] = l
			f.cnt[i] = c
		} else if l == f.len[i] {
			f.cnt[i] += c
		}
		i += i & -i
	}
}

func (f *Fenwick) query(i int) (int, float64) {
	bestL := 0
	bestC := 0.0
	for i > 0 {
		if f.len[i] > bestL {
			bestL = f.len[i]
			bestC = f.cnt[i]
		} else if f.len[i] == bestL {
			bestC += f.cnt[i]
		}
		i -= i & -i
	}
	return bestL, bestC
}

// unique returns sorted unique slice of a.
func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	res := []int{a[0]}
	for _, v := range a[1:] {
		if v != res[len(res)-1] {
			res = append(res, v)
		}
	}
	return res
}

// minReplacements computes the minimal number of replacements
// to make the array almost increasing.
func minReplacements(a []int) int {
	n := len(a)
	vals := append([]int(nil), a...)
	sort.Ints(vals)
	vals = unique(vals)
	m := len(vals)
	idx := make(map[int]int, m)
	for i, v := range vals {
		idx[v] = i + 1
	}

	L := make([]int, n)
	C1 := make([]float64, n)
	fw := NewFenwick(m + 2)
	for i, x := range a {
		pos := idx[x]
		l, c := fw.query(pos - 1)
		if l == 0 {
			c = 1
		}
		L[i] = l + 1
		C1[i] = c
		fw.update(pos, l+1, c)
	}

	LIS := 0
	for _, v := range L {
		if v > LIS {
			LIS = v
		}
	}
	total := 0.0
	for i, v := range L {
		if v == LIS {
			total += C1[i]
		}
	}

	R := make([]int, n)
	C2 := make([]float64, n)
	fw2 := NewFenwick(m + 2)
	for i := n - 1; i >= 0; i-- {
		pos := m - idx[a[i]] + 1
		l, c := fw2.query(pos - 1)
		if l == 0 {
			c = 1
		}
		R[i] = l + 1
		C2[i] = c
		fw2.update(pos, l+1, c)
	}

	nonessential := false
	for i := 0; i < n; i++ {
		if L[i]+R[i]-1 == LIS {
			if C1[i]*C2[i] < total-1e-9 {
				nonessential = true
				break
			}
		} else {
			nonessential = true
			break
		}
	}
	if nonessential {
		return n - (LIS + 1)
	}
	return n - LIS
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	fmt.Println(minReplacements(a))
}
