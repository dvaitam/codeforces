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
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.tree[i] += v
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

func (f *Fenwick) RangeSum(l, r int) int {
	if l <= r {
		return f.Sum(r) - f.Sum(l-1)
	}
	return f.Sum(f.n) - f.Sum(l-1) + f.Sum(r)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	posMap := make(map[int][]int)
	for i := 1; i <= n; i++ {
		posMap[a[i]] = append(posMap[a[i]], i)
	}
	vals := make([]int, 0, len(posMap))
	for v := range posMap {
		vals = append(vals, v)
	}
	sort.Ints(vals)

	fw := NewFenwick(n)
	for i := 1; i <= n; i++ {
		fw.Add(i, 1)
	}
	cur := 1
	var ops int64
	for _, v := range vals {
		idxs := posMap[v]
		start := sort.Search(len(idxs), func(i int) bool { return idxs[i] >= cur })
		for i := start; i < len(idxs); i++ {
			idx := idxs[i]
			ops += int64(fw.RangeSum(cur, idx))
			fw.Add(idx, -1)
			cur = idx
		}
		for i := 0; i < start; i++ {
			idx := idxs[i]
			ops += int64(fw.RangeSum(cur, idx))
			fw.Add(idx, -1)
			cur = idx
		}
	}

	fmt.Fprintln(writer, ops)
}
