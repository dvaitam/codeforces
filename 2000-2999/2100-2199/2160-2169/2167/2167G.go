package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

type fenwick struct {
	tree []int64
}

func newFenwick(size int) *fenwick {
	return &fenwick{tree: make([]int64, size)}
}

func (f *fenwick) update(idx int, val int64) {
	for idx > 0 && idx < len(f.tree) {
		if val > f.tree[idx] {
			f.tree[idx] = val
		}
		idx += idx & -idx
	}
}

func (f *fenwick) query(idx int) int64 {
	var res int64
	for idx > 0 {
		if f.tree[idx] > res {
			res = f.tree[idx]
		}
		idx -= idx & -idx
	}
	return res
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	answers := make([]int64, t)

	for tc := 0; tc < t; tc++ {
		n := fs.nextInt()
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
		}
		cost := make([]int64, n)
		for i := 0; i < n; i++ {
			cost[i] = int64(fs.nextInt())
		}

		vals := make([]int, n)
		copy(vals, a)
		sort.Ints(vals)
		uniq := make([]int, 0, n)
		for _, v := range vals {
			if len(uniq) == 0 || uniq[len(uniq)-1] != v {
				uniq = append(uniq, v)
			}
		}

		comp := make(map[int]int, len(uniq))
		for i, v := range uniq {
			comp[v] = i + 1
		}

		bit := newFenwick(len(uniq) + 2)
		var total, bestKeep int64
		for i := 0; i < n; i++ {
			total += cost[i]
			idx := comp[a[i]]
			cur := bit.query(idx) + cost[i]
			if cur > bestKeep {
				bestKeep = cur
			}
			bit.update(idx, cur)
		}
		answers[tc] = total - bestKeep
	}

	for _, ans := range answers {
		fmt.Fprintln(out, ans)
	}
}
