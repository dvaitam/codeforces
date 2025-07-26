package main

import (
	"bufio"
	"fmt"
	"os"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT { b := &BIT{n: n, tree: make([]int, n+2)}; return b }
func (b *BIT) Add(i, val int) {
	for i <= b.n {
		b.tree[i] += val
		i += i & -i
	}
}
func (b *BIT) Sum(i int) int {
	res := 0
	for i > 0 {
		res += b.tree[i]
		i -= i & -i
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	var s string
	fmt.Fscan(in, &s)
	parent := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] == x {
			return x
		}
		parent[x] = find(parent[x])
		return parent[x]
	}
	order := make([]int, 0, n)
	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		x := find(l)
		for x <= r {
			order = append(order, x)
			parent[x] = x + 1
			x = find(x)
		}
	}
	pos := make([]int, n+1)
	for i := range pos {
		pos[i] = -1
	}
	for idx, val := range order {
		pos[val] = idx + 1
	}
	bit := NewBIT(len(order))
	b := []byte(s)
	ones := 0
	for i := 1; i <= n; i++ {
		if b[i-1] == '1' {
			ones++
		}
		if p := pos[i]; p > 0 {
			bit.Add(p, int(b[i-1]-'0'))
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for qi := 0; qi < q; qi++ {
		var x int
		fmt.Fscan(in, &x)
		if b[x-1] == '1' {
			b[x-1] = '0'
			ones--
			if p := pos[x]; p > 0 {
				bit.Add(p, -1)
			}
		} else {
			b[x-1] = '1'
			ones++
			if p := pos[x]; p > 0 {
				bit.Add(p, 1)
			}
		}
		k := ones
		if k > len(order) {
			k = len(order)
		}
		onesIn := bit.Sum(k)
		if qi+1 == q {
			fmt.Fprintf(out, "%d\n", k-onesIn)
		} else {
			fmt.Fprintf(out, "%d ", k-onesIn)
		}
	}
	out.Flush()
}
