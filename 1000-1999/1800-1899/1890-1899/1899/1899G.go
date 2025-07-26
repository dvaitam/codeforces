package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{make([]int, n+2)}
}

func (f *Fenwick) Add(i, delta int) {
	for i < len(f.tree) {
		f.tree[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	res := 0
	for i > 0 {
		res += f.tree[i]
		i -= i & -i
	}
	return res
}

func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
		}

		tin := make([]int, n+1)
		tout := make([]int, n+1)
		type item struct{ v, p, state int }
		stack := []item{{1, 0, 0}}
		timer := 0
		for len(stack) > 0 {
			it := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if it.state == 0 {
				timer++
				tin[it.v] = timer
				stack = append(stack, item{it.v, it.p, 1})
				for i := len(g[it.v]) - 1; i >= 0; i-- {
					u := g[it.v][i]
					if u == it.p {
						continue
					}
					stack = append(stack, item{u, it.v, 0})
				}
			} else {
				tout[it.v] = timer
			}
		}

		type Query struct{ l, r, x int }
		queries := make([]Query, q)
		qByR := make([][]int, n+1)
		qByL := make([][]int, n+1)
		for i := 0; i < q; i++ {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			queries[i] = Query{l, r, x}
			qByR[r] = append(qByR[r], i)
			qByL[l-1] = append(qByL[l-1], i)
		}

		fen := NewFenwick(n)
		resR := make([]int, q)
		resL := make([]int, q)

		for _, id := range qByR[0] {
			x := queries[id].x
			resR[id] = fen.RangeSum(tin[x], tout[x])
		}
		for _, id := range qByL[0] {
			x := queries[id].x
			resL[id] = fen.RangeSum(tin[x], tout[x])
		}

		for i := 1; i <= n; i++ {
			fen.Add(tin[p[i]], 1)
			for _, id := range qByR[i] {
				x := queries[id].x
				resR[id] = fen.RangeSum(tin[x], tout[x])
			}
			for _, id := range qByL[i] {
				x := queries[id].x
				resL[id] = fen.RangeSum(tin[x], tout[x])
			}
		}

		for i := 0; i < q; i++ {
			if resR[i]-resL[i] > 0 {
				fmt.Fprintln(writer, "Yes")
			} else {
				fmt.Fprintln(writer, "No")
			}
		}
	}
}
