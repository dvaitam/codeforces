package main

import (
	"bufio"
	"fmt"
	"os"
)

// Fenwick tree supporting prefix maximum queries.
type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) Update(idx int, val int64) {
	for idx <= f.n {
		if val > f.bit[idx] {
			f.bit[idx] = val
		}
		idx += idx & -idx
	}
}

func (f *Fenwick) Query(idx int) int64 {
	if idx <= 0 {
		return 0
	}
	if idx > f.n {
		idx = f.n
	}
	res := int64(0)
	for idx > 0 {
		if f.bit[idx] > res {
			res = f.bit[idx]
		}
		idx &= idx - 1
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

	attacks := make([][]struct {
		v   int
		val int64
	}, n+1)

	for i := 1; i <= n; i++ {
		var k int
		fmt.Fscan(in, &k)
		damage := make([]int64, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &damage[j])
		}
		for j := 0; j < k; j++ {
			var b int
			fmt.Fscan(in, &b)
			d := m - b
			if d < 0 || d > m {
				continue
			}
			if i-1 < d {
				// not enough previous warriors to destroy d barricades
				continue
			}
			v := i - d
			attacks[i] = append(attacks[i], struct {
				v   int
				val int64
			}{v, damage[j]})
		}
	}

	ft := NewFenwick(n + 2)
	ft.Update(1, 0) // starting state: v=0
	ans := int64(0)

	for i := 1; i <= n; i++ {
		// evaluate all attacks of warrior i using values from previous steps
		tmp := make(map[int]int64)
		for _, at := range attacks[i] {
			cand := ft.Query(at.v) + at.val
			if cand > tmp[at.v] {
				tmp[at.v] = cand
			}
		}
		// apply best results for this warrior
		for v, val := range tmp {
			if val > ans {
				ans = val
			}
			ft.Update(v+1, val)
		}
	}

	fmt.Fprintln(out, ans)
}
