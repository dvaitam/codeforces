package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// BIT for prefix max queries and point updates
type BIT struct {
	n    int
	tree []int64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) Update(i int, val int64) {
	for ; i <= b.n; i += i & -i {
		if b.tree[i] < val {
			b.tree[i] = val
		}
	}
}

func (b *BIT) Query(i int) int64 {
	res := int64(math.MinInt64 / 4)
	for ; i > 0; i -= i & -i {
		if b.tree[i] > res {
			res = b.tree[i]
		}
	}
	return res
}

func unique64(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	const negInf = int64(math.MinInt64 / 4)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + arr[i]
		}
		vals := append([]int64(nil), pref...)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		vals = unique64(vals)
		m := len(vals)
		comp := make(map[int64]int, m)
		for i, v := range vals {
			comp[v] = i + 1
		}
		idx := make([]int, n+1)
		for i := 0; i <= n; i++ {
			idx[i] = comp[pref[i]]
		}
		eq := make([]int64, m+1)
		for i := range eq {
			eq[i] = negInf
		}
		bitLess := NewBIT(m + 2)
		bitGreater := NewBIT(m + 2)

		dp := int64(0)
		id0 := idx[0]
		eq[id0] = 0
		bitLess.Update(id0, 0)
		bitGreater.Update(m-id0+1, 0)

		for i := 1; i <= n; i++ {
			id := idx[i]
			best := eq[id]
			v := bitLess.Query(id - 1)
			if v != negInf {
				if cand := v + int64(i); cand > best {
					best = cand
				}
			}
			w := bitGreater.Query(m - id)
			if w != negInf {
				if cand := w - int64(i); cand > best {
					best = cand
				}
			}
			dp = best
			if dp > eq[id] {
				eq[id] = dp
			}
			bitLess.Update(id, dp-int64(i))
			bitGreater.Update(m-id+1, dp+int64(i))
		}
		fmt.Fprintln(out, dp)
	}
}
