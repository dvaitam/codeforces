package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 1 << 60

type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) add(idx int, val int64) {
	for idx <= f.n+1 {
		f.tree[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) rangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	f.add(l+1, val)
	f.add(r+2, -val)
}

func (f *Fenwick) sum(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func (f *Fenwick) pointQuery(i int) int64 {
	return f.sum(i + 1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}
	p := make([]int, n)
	for i := range p {
		fmt.Fscan(reader, &p[i])
	}
	var m int
	fmt.Fscan(reader, &m)
	b := make([]int, m)
	for i := range b {
		fmt.Fscan(reader, &b[i])
	}

	// dp array indices 0..m
	dpBase := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		dpBase[i] = INF
	}

	bit := NewFenwick(m + 1)

	// helper functions
	query := func(idx int) int64 {
		return dpBase[idx] + bit.pointQuery(idx)
	}
	setVal := func(idx int, val int64) {
		dpBase[idx] = val - bit.pointQuery(idx)
	}
	rangeAdd := func(l, r int, val int64) {
		bit.rangeAdd(l, r, val)
	}

	for i := 0; i < n; i++ {
		x := a[i]
		cost := int64(p[i])
		pos := sort.Search(len(b), func(j int) bool { return b[j] >= x })
		if pos == m { // greater than all b
			rangeAdd(0, m, cost)
			continue
		}
		old := query(pos)
		if pos > 0 {
			rangeAdd(0, pos-1, cost)
		}
		addVal := cost
		if addVal > 0 {
			addVal = 0
		}
		rangeAdd(pos, m, addVal)
		if b[pos] == x {
			cur := query(pos + 1)
			if old < cur {
				setVal(pos+1, old)
			}
		}
	}

	ans := query(m)
	if ans >= INF/2 {
		fmt.Fprintln(writer, "NO")
	} else {
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, ans)
	}
}
