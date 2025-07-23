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

func (b *BIT) Update(i, val int) {
	for i <= b.n {
		if val > b.tree[i] {
			b.tree[i] = val
		}
		i += i & -i
	}
}

func (b *BIT) Query(i int) int {
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	h := make([]int, n)
	values := make([]int, n+m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
		values[i] = h[i]
	}

	a := make([]int, m)
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &a[i], &b[i])
		a[i]--
		values[n+i] = b[i]
	}

	sort.Ints(values)
	uniq := values[:0]
	for _, v := range values {
		if len(uniq) == 0 || v != uniq[len(uniq)-1] {
			uniq = append(uniq, v)
		}
	}
	idx := make(map[int]int, len(uniq))
	for i, v := range uniq {
		idx[v] = i + 1
	}
	k := len(uniq)

	queriesAt := make([][]int, n)
	for i := 0; i < m; i++ {
		queriesAt[a[i]] = append(queriesAt[a[i]], i)
	}

	L := make([]int, n)
	prefix := make([]int, m)
	bit := NewBIT(k)
	for i := 0; i < n; i++ {
		for _, qi := range queriesAt[i] {
			pos := idx[b[qi]]
			prefix[qi] = bit.Query(pos-1) + 1
		}
		pos := idx[h[i]]
		L[i] = bit.Query(pos-1) + 1
		bit.Update(pos, L[i])
	}

	R := make([]int, n)
	suffix := make([]int, m)
	bit2 := NewBIT(k)
	for i := n - 1; i >= 0; i-- {
		for _, qi := range queriesAt[i] {
			pos := idx[b[qi]]
			rev := k - pos + 1
			suffix[qi] = bit2.Query(rev-1) + 1
		}
		pos := idx[h[i]]
		rev := k - pos + 1
		R[i] = bit2.Query(rev-1) + 1
		bit2.Update(rev, R[i])
	}

	lis := 0
	for i := 0; i < n; i++ {
		if L[i]+R[i]-1 > lis {
			lis = L[i] + R[i] - 1
		}
	}
	cnt := make(map[int]int)
	for i := 0; i < n; i++ {
		if L[i]+R[i]-1 == lis {
			cnt[L[i]]++
		}
	}
	critical := make([]bool, n)
	for i := 0; i < n; i++ {
		if L[i]+R[i]-1 == lis && cnt[L[i]] == 1 {
			critical[i] = true
		}
	}

	for i := 0; i < m; i++ {
		base := lis
		if critical[a[i]] {
			base = lis - 1
		}
		cand := prefix[i] + suffix[i] - 1
		if cand < base {
			fmt.Fprintln(out, base)
		} else {
			fmt.Fprintln(out, cand)
		}
	}
}
