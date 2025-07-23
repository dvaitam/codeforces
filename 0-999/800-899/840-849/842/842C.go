package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func divisors(x int) []int {
	res := []int{}
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			res = append(res, i)
			if i*i != x {
				res = append(res, x/i)
			}
		}
	}
	sort.Ints(res)
	return res
}

type item struct {
	v, p, depth int
	g1, g0      int
	state       int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	rootVal := a[1]
	divs := divisors(rootVal)
	idx := make(map[int]int, len(divs))
	for i, d := range divs {
		idx[d] = i
	}

	divIdx := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		for j, d := range divs {
			if a[i]%d == 0 {
				divIdx[i] = append(divIdx[i], j)
			}
		}
	}

	cnt := make([]int, len(divs))
	ans := make([]int, n+1)

	stack := []item{{v: 1, p: 0, depth: 1, g1: 0, g0: 0, state: 0}}
	for len(stack) > 0 {
		it := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if it.state == 0 {
			g1 := gcd(it.g1, a[it.v])
			var g0 int
			if it.p == 0 {
				g0 = 0
			} else {
				g0 = gcd(it.g0, a[it.v])
			}
			stack = append(stack, item{v: it.v, p: it.p, depth: it.depth, g1: g1, g0: g0, state: 1})
			for _, id := range divIdx[it.v] {
				cnt[id]++
			}
			best := g1
			if g0 > best {
				best = g0
			}
			for i, d := range divs {
				if cnt[i] >= it.depth-1 && d > best {
					best = d
				}
			}
			ans[it.v] = best
			for _, to := range adj[it.v] {
				if to == it.p {
					continue
				}
				stack = append(stack, item{v: to, p: it.v, depth: it.depth + 1, g1: g1, g0: g0, state: 0})
			}
		} else {
			for _, id := range divIdx[it.v] {
				cnt[id]--
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
