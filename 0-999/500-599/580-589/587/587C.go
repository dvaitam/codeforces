package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const LOG = 17
const LIM = 10

func merge(a, b []int) []int {
	res := make([]int, 0, LIM)
	i, j := 0, 0
	for len(res) < LIM && (i < len(a) || j < len(b)) {
		var x int
		if j >= len(b) || (i < len(a) && a[i] < b[j]) {
			x = a[i]
			i++
		} else {
			x = b[j]
			j++
		}
		if len(res) == 0 || res[len(res)-1] != x {
			res = append(res, x)
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	people := make([][]int, n+1)
	for i := 1; i <= m; i++ {
		var c int
		fmt.Fscan(reader, &c)
		people[c] = append(people[c], i)
	}
	for i := 1; i <= n; i++ {
		sort.Ints(people[i])
		if len(people[i]) > LIM {
			people[i] = people[i][:LIM]
		}
	}

	parent := make([][]int, LOG+1)
	upvals := make([][][]int, LOG+1)
	for i := 0; i <= LOG; i++ {
		parent[i] = make([]int, n+1)
		upvals[i] = make([][]int, n+1)
	}

	depth := make([]int, n+1)
	qarr := make([]int, 0, n)
	qarr = append(qarr, 1)
	parent[0][1] = 0
	depth[1] = 0
	for idx := 0; idx < len(qarr); idx++ {
		v := qarr[idx]
		for _, to := range adj[v] {
			if to == parent[0][v] {
				continue
			}
			parent[0][to] = v
			depth[to] = depth[v] + 1
			qarr = append(qarr, to)
		}
	}

	for v := 1; v <= n; v++ {
		p := parent[0][v]
		if p == 0 {
			upvals[0][v] = append([]int(nil), people[v]...)
		} else {
			upvals[0][v] = merge(people[v], people[p])
		}
	}

	for k := 1; k <= LOG; k++ {
		for v := 1; v <= n; v++ {
			parent[k][v] = parent[k-1][parent[k-1][v]]
			upvals[k][v] = merge(upvals[k-1][v], upvals[k-1][parent[k-1][v]])
		}
	}

	for ; q > 0; q-- {
		var v, u, a int
		fmt.Fscan(reader, &v, &u, &a)
		res := make([]int, 0)
		if depth[v] < depth[u] {
			v, u = u, v
		}
		diff := depth[v] - depth[u]
		for k := LOG; k >= 0; k-- {
			if diff>>k&1 == 1 {
				res = merge(res, upvals[k][v])
				v = parent[k][v]
			}
		}
		if v == u {
			res = merge(res, people[v])
		} else {
			for k := LOG; k >= 0; k-- {
				if parent[k][v] != parent[k][u] {
					res = merge(res, upvals[k][v])
					res = merge(res, upvals[k][u])
					v = parent[k][v]
					u = parent[k][u]
				}
			}
			res = merge(res, upvals[0][v])
			res = merge(res, upvals[0][u])
		}
		if len(res) > a {
			res = res[:a]
		}
		fmt.Fprint(writer, len(res))
		for _, x := range res {
			fmt.Fprint(writer, " ", x)
		}
		fmt.Fprintln(writer)
	}
}
