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
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) Reset() {
	for i := range f.tree {
		f.tree[i] = 0
	}
}

func (f *Fenwick) Add(i int, val int64) {
	for i <= f.n {
		f.tree[i] += val
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

func query(cnt, sum *Fenwick, weights []int64, limit int64) int {
	if limit <= 0 {
		return 0
	}
	n := sum.n
	idx := 0
	s := int64(0)
	c := int64(0)
	bit := 1
	for bit<<1 <= n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= n && s+sum.tree[next] <= limit {
			s += sum.tree[next]
			c += cnt.tree[next]
			idx = next
		}
		bit >>= 1
	}
	remain := limit - s
	if idx < n {
		weight := weights[idx+1]
		avail := cnt.Sum(idx+1) - cnt.Sum(idx)
		more := remain / weight
		if more > avail {
			more = avail
		}
		c += more
	}
	return int(c)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var L int64
		fmt.Fscan(in, &n, &L)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i], &b[i])
		}
		msgs := make([]struct{ a, b int64 }, n)
		for i := 0; i < n; i++ {
			msgs[i].a = a[i]
			msgs[i].b = b[i]
		}
		sort.Slice(msgs, func(i, j int) bool { return msgs[i].b < msgs[j].b })
		aSorted := make([]int64, n)
		for i := 0; i < n; i++ {
			aSorted[i] = msgs[i].a
		}
		unique := append([]int64(nil), aSorted...)
		sort.Slice(unique, func(i, j int) bool { return unique[i] < unique[j] })
		m := 0
		for i := 0; i < len(unique); i++ {
			if i == 0 || unique[i] != unique[m-1] {
				unique[m] = unique[i]
				m++
			}
		}
		unique = unique[:m]
		idxMap := make(map[int64]int, m)
		for i, v := range unique {
			idxMap[v] = i + 1
		}
		// weights array is 1-indexed for Fenwick queries
		weights := make([]int64, m+1)
		for i := 0; i < m; i++ {
			weights[i+1] = unique[i]
		}
		cntBIT := NewFenwick(m)
		sumBIT := NewFenwick(m)
		ans := 0
		for i := 0; i < n; i++ {
			if msgs[i].a <= L {
				if ans < 1 {
					ans = 1
				}
			}
			cntBIT.Reset()
			sumBIT.Reset()
			for j := i + 1; j < n; j++ {
				// base cost using i and j as extremes
				cost := msgs[j].b - msgs[i].b + msgs[i].a + msgs[j].a
				if cost <= L {
					leftover := L - cost
					c := query(cntBIT, sumBIT, weights, leftover)
					if c+2 > ans {
						ans = c + 2
					}
				}
				// add j to interior for next iterations
				idx := idxMap[msgs[j].a]
				cntBIT.Add(idx, 1)
				sumBIT.Add(idx, msgs[j].a)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
