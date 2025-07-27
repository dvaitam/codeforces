package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+2)}
}

func (f *Fenwick) Add(i, delta int) {
	for i < len(f.tree) {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	res := 0
	for i > 0 {
		res += f.tree[i]
		i -= i & -i
	}
	return res
}

func (f *Fenwick) FindByOrder(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= len(f.tree) {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next < len(f.tree) && f.tree[next] < k {
			idx = next
			k -= f.tree[next]
		}
		bit >>= 1
	}
	return idx + 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	size := n + q + 5
	ft := NewFenwick(size)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		ft.Add(x, 1)
	}
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(reader, &k)
		if k > 0 {
			ft.Add(k, 1)
		} else {
			idx := ft.FindByOrder(-k)
			ft.Add(idx, -1)
		}
	}
	if ft.Sum(size-1) == 0 {
		fmt.Fprint(writer, 0)
		return
	}
	ans := ft.FindByOrder(1)
	fmt.Fprint(writer, ans)
}
