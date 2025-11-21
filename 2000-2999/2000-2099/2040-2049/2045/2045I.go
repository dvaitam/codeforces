package main

import (
	"bufio"
	"fmt"
	"os"
)

type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) add(idx int, delta int64) {
	for idx <= b.n {
		b.tree[idx] += delta
		idx += idx & -idx
	}
}

func (b *BIT) sum(idx int) int64 {
	var res int64
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (b *BIT) rangeSum(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.sum(r) - b.sum(l-1)
}

type Query struct {
	l   int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	pos := make([][]int, m+1)
	for i := 1; i <= n; i++ {
		v := a[i]
		pos[v] = append(pos[v], i)
	}

	P := 0
	for v := 1; v <= m; v++ {
		if len(pos[v]) > 0 {
			P++
		}
	}

	queriesByR := make([][]Query, n+1)
	var results []int64

	for v := 1; v <= m; v++ {
		occ := pos[v]
		for i := 0; i+1 < len(occ); i++ {
			l := occ[i] + 1
			r := occ[i+1] - 1
			if l <= r {
				idx := len(results)
				results = append(results, 0)
				queriesByR[r] = append(queriesByR[r], Query{l: l, idx: idx})
			}
		}
	}

	bit := NewBIT(n)
	lastOcc := make([]int, m+1)

	for i := 1; i <= n; i++ {
		v := a[i]
		if lastOcc[v] != 0 {
			bit.add(lastOcc[v], -1)
		}
		bit.add(i, 1)
		lastOcc[v] = i
		for _, q := range queriesByR[i] {
			results[q.idx] = bit.rangeSum(q.l, i)
		}
	}

	var totalC int64
	for _, val := range results {
		totalC += val
	}

	ans := int64(P)*(int64(m)-1) + totalC
	fmt.Fprintln(out, ans)
}
