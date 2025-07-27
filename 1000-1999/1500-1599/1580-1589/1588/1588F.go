package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) Add(i int, delta int64) {
	for i <= f.n {
		f.bit[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int64 {
	var res int64
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

func (f *Fenwick) RangeSum(l, r int) int64 {
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)

	a := make([]int64, n+1)
	fw := NewFenwick(n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		fw.Add(i, a[i])
	}

	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
	}

	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	visited := make([]bool, n+1)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		switch t {
		case 1:
			var l, r int
			fmt.Fscan(in, &l, &r)
			fmt.Fprintln(out, fw.RangeSum(l, r))
		case 2:
			var v int
			var x int64
			fmt.Fscan(in, &v, &x)
			for i := 1; i <= n; i++ {
				visited[i] = false
			}
			cur := v
			for !visited[cur] {
				visited[cur] = true
				fw.Add(cur, x)
				cur = p[cur]
			}
		case 3:
			var i, j int
			fmt.Fscan(in, &i, &j)
			p[i], p[j] = p[j], p[i]
		}
	}
}
