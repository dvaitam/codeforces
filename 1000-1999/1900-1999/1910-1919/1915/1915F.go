package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+1)}
}

func (f *Fenwick) Add(i int, v int64) {
	for i <= f.n {
		f.tree[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += f.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		people := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &people[i][0], &people[i][1])
		}
		sort.Slice(people, func(i, j int) bool { return people[i][0] < people[j][0] })
		bs := make([]int, n)
		for i := 0; i < n; i++ {
			bs[i] = people[i][1]
		}
		sortedBs := append([]int(nil), bs...)
		sort.Ints(sortedBs)
		comp := make(map[int]int, n)
		for i, v := range sortedBs {
			comp[v] = i + 1
		}
		bit := NewFenwick(n)
		ans := int64(0)
		for i := 0; i < n; i++ {
			idx := comp[bs[i]]
			ans += int64(i) - bit.Sum(idx)
			bit.Add(idx, 1)
		}
		fmt.Fprintln(out, ans)
	}
}
