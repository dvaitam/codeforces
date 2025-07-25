package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(idx, val int) {
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	s := 0
	for idx > 0 {
		s += b.tree[idx]
		idx -= idx & -idx
	}
	return s
}

func (b *BIT) FindKth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}

type Elem struct {
	val int
	idx int
}

type Query struct {
	k   int
	pos int
	id  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	elems := make([]Elem, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		elems[i] = Elem{val: arr[i], idx: i + 1}
	}

	sort.Slice(elems, func(i, j int) bool {
		if elems[i].val == elems[j].val {
			return elems[i].idx < elems[j].idx
		}
		return elems[i].val > elems[j].val
	})

	var m int
	fmt.Fscan(in, &m)
	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &queries[i].k, &queries[i].pos)
		queries[i].id = i
	}

	sort.Slice(queries, func(i, j int) bool { return queries[i].k < queries[j].k })

	bit := NewBIT(n)
	answers := make([]int, m)
	inserted := 0
	for _, q := range queries {
		for inserted < q.k {
			bit.Add(elems[inserted].idx, 1)
			inserted++
		}
		idx := bit.FindKth(q.pos)
		answers[q.id] = arr[idx-1]
	}

	for i := 0; i < m; i++ {
		fmt.Fprintln(out, answers[i])
	}
}
