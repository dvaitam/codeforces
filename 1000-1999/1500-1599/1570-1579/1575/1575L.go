package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// BIT implements a Fenwick tree for maximum queries.
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Update(i, val int) {
	i++
	for i <= b.n+1 {
		if val > b.tree[i] {
			b.tree[i] = val
		}
		i += i & -i
	}
}

func (b *BIT) Query(i int) int {
	if i < 0 {
		return 0
	}
	i++
	res := 0
	for i > 0 {
		if b.tree[i] > res {
			res = b.tree[i]
		}
		i -= i & -i
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	groups := make(map[int][]int)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] <= i {
			groups[a[i]] = append(groups[a[i]], i)
		}
	}
	values := make([]int, 0, len(groups))
	for v := range groups {
		values = append(values, v)
	}
	sort.Ints(values)
	bit := NewBIT(n)
	ans := 0
	for _, v := range values {
		idxs := groups[v]
		updates := make([][2]int, len(idxs))
		for j, idx := range idxs {
			d := idx - v
			best := bit.Query(d)
			updates[j] = [2]int{d, best + 1}
			if best+1 > ans {
				ans = best + 1
			}
		}
		for _, u := range updates {
			bit.Update(u[0], u[1])
		}
	}
	fmt.Println(ans)
}
