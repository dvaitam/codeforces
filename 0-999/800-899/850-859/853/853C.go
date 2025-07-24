package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type query struct {
	l, d, r, u int
}

type item struct {
	x   int
	y   int
	idx int
	id  int
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, v int) {
	for i <= b.n {
		b.tree[i] += v
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

func choose2(x int64) int64 {
	if x < 2 {
		return 0
	}
	return x * (x - 1) / 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}
	queries := make([]query, q)
	items := make([]item, 0, 4*q)
	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		queries[i] = query{l, d, r, u}
		items = append(items, item{x: l - 1, y: d - 1, idx: i, id: 0})
		items = append(items, item{x: l - 1, y: u, idx: i, id: 1})
		items = append(items, item{x: r, y: d - 1, idx: i, id: 2})
		items = append(items, item{x: r, y: u, idx: i, id: 3})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].x < items[j].x })

	bit := NewBIT(n)
	res := make([][4]int, q)
	nextCol := 1
	for _, it := range items {
		for nextCol <= it.x {
			bit.Add(p[nextCol-1], 1)
			nextCol++
		}
		val := 0
		if it.y >= 1 {
			val = bit.Sum(it.y)
		} else if it.y == 0 {
			val = bit.Sum(0)
		} else {
			val = 0
		}
		res[it.idx][it.id] = val
	}

	total := int64(n) * int64(n-1) / 2
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for i, qy := range queries {
		a := res[i][0]
		b := res[i][1]
		c := res[i][2]
		d2 := res[i][3]
		L := qy.l - 1
		R := n - qy.r
		Dcnt := qy.d - 1
		Ucnt := n - qy.u
		LD := a
		LU := L - b
		RD := Dcnt - c
		RU := Ucnt - (qy.r - d2)
		non := choose2(int64(L)) + choose2(int64(R)) + choose2(int64(Dcnt)) + choose2(int64(Ucnt)) -
			choose2(int64(LD)) - choose2(int64(LU)) - choose2(int64(RD)) - choose2(int64(RU))
		ans := total - non
		fmt.Fprintln(out, ans)
	}
}
