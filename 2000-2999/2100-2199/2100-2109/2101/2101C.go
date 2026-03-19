package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// DSU for slot scheduling
type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := range p {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

type Item struct {
	pos      int
	deadline int
	weight   int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		items := make([]Item, n)
		for i := 0; i < n; i++ {
			d := 2 * a[i]
			if d > 2*n {
				d = 2 * n
			}
			items[i] = Item{
				pos:      i + 1,
				deadline: d,
				weight:   abs(2*(i+1) - n - 1),
			}
		}

		sort.Slice(items, func(i, j int) bool {
			if items[i].weight == items[j].weight {
				return items[i].pos > items[j].pos
			}
			return items[i].weight > items[j].weight
		})

		dsu := NewDSU(2 * n)
		accepted := make([]int, 0, n)

		for _, item := range items {
			slot := dsu.Find(item.deadline)
			if slot > 0 {
				accepted = append(accepted, item.pos)
				dsu.parent[slot] = slot - 1
			}
		}

		sort.Ints(accepted)
		m := len(accepted)
		k := m / 2

		var total int64
		for i := 0; i < k; i++ {
			total -= int64(accepted[i])
		}
		for i := m - k; i < m; i++ {
			total += int64(accepted[i])
		}

		fmt.Fprintln(out, total)
	}
}
