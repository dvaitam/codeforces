package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	first  int
	second int
}

func merge(a, b pair) pair {
	// keep two largest distinct values
	x1, x2 := a.first, a.second
	candidates := []int{b.first, b.second}
	for _, v := range candidates {
		if v == x1 || v == x2 || v == -1 {
			continue
		}
		if v > x1 {
			x2 = x1
			x1 = v
		} else if v > x2 {
			x2 = v
		}
	}
	return pair{x1, x2}
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
		w := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &w[i])
		}
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		tin := make([]int, n)
		tout := make([]int, n)
		order := make([]int, 0, n)

		type item struct {
			v        int
			parent   int
			processed bool
		}
		st := []item{{0, -1, false}}
		for len(st) > 0 {
			cur := st[len(st)-1]
			st = st[:len(st)-1]
			if cur.processed {
				tout[cur.v] = len(order)
				continue
			}
			tin[cur.v] = len(order)
			order = append(order, cur.v)
			st = append(st, item{cur.v, cur.parent, true})
			for _, to := range g[cur.v] {
				if to == cur.parent {
					continue
				}
				st = append(st, item{to, cur.v, false})
			}
		}

		// build arrays of node values in Euler order
		vals := make([]int, n)
		for idx, v := range order {
			vals[idx] = w[v]
		}

		pref := make([]pair, n)
		for i := 0; i < n; i++ {
			if i == 0 {
				pref[i] = pair{vals[i], -1}
			} else {
				pref[i] = merge(pref[i-1], pair{vals[i], -1})
			}
		}

		suff := make([]pair, n)
		for i := n - 1; i >= 0; i-- {
			if i == n-1 {
				suff[i] = pair{vals[i], -1}
			} else {
				suff[i] = merge(suff[i+1], pair{vals[i], -1})
			}
		}

		answer := 0
		for v := 0; v < n && answer == 0; v++ {
			l, r := tin[v], tout[v]
			outside := pair{-1, -1}
			if l > 0 {
				outside = merge(outside, pref[l-1])
			}
			if r < n {
				outside = merge(outside, suff[r])
			}
			top1, top2 := outside.first, outside.second
			if top1 == -1 {
				continue
			}
			if top1 > w[v] && (top2 == -1 || top2 <= w[v]) {
				answer = v + 1
			}
		}

		fmt.Fprintln(out, answer)
	}
}
