package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	l   int
	r   int
	idx int
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	segs := make([]Segment, n)
	rights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &segs[i].l, &segs[i].r)
		segs[i].idx = i
		rights[i] = segs[i].r
	}

	// Compress right endpoints
	sortedRights := make([]int, n)
	copy(sortedRights, rights)
	sort.Ints(sortedRights)
	mp := make(map[int]int, n)
	for i, v := range sortedRights {
		mp[v] = i + 1
	}
	for i := range segs {
		segs[i].r = mp[segs[i].r]
	}

	sort.Slice(segs, func(i, j int) bool {
		return segs[i].l < segs[j].l
	})

	bit := NewBIT(n)
	ans := make([]int, n)

	for i := n - 1; i >= 0; i-- {
		rp := segs[i].r
		ans[segs[i].idx] = bit.Sum(rp)
		bit.Add(rp, 1)
	}

	for i := 0; i < n; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
