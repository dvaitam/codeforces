package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	f := &Fenwick{n: n, tree: make([]int, n+2)}
	return f
}

func (f *Fenwick) Add(i, delta int) {
	for i <= f.n {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func (f *Fenwick) Kth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= f.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= f.n && f.tree[next] < k {
			k -= f.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

func permutationValue(b []int) []int {
	n := len(b)
	fw := NewFenwick(n)
	for i := 1; i <= n; i++ {
		fw.Add(i, 1)
	}
	p := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		k := b[i] + 1
		x := fw.Kth(k)
		p[i] = x
		fw.Add(x, -1)
	}
	return p
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var i, x int
			fmt.Fscan(reader, &i, &x)
			b[i-1] = x
		} else {
			var i int
			fmt.Fscan(reader, &i)
			p := permutationValue(b)
			fmt.Fprintln(writer, p[i-1])
		}
	}
}
