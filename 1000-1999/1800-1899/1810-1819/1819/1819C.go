package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, rt, pos int
	to         [][]int
	deg        []int
	res        []int
	out        = bufio.NewWriter(os.Stdout)
)

func failure() {
	fmt.Fprintln(os.Stdout, "No")
	os.Exit(0)
}

func dfs(u, fa int, chs bool) {
	d := 0
	for _, v := range to[u] {
		if v != fa && deg[v] > 1 {
			d++
		}
	}
	sub := d
	if u == rt {
		sub = d - 1
	}
	if sub > 1 {
		failure()
	}
	if u == rt {
		var x int
		pos = 1
		res[pos] = rt
		for _, v := range to[u] {
			if deg[v] > 1 {
				x = v
				dfs(v, u, false)
				break
			}
		}
		for _, v := range to[u] {
			if deg[v] == 1 {
				pos++
				res[pos] = v
			}
		}
		for _, v := range to[u] {
			if deg[v] > 1 && v != x {
				dfs(v, u, true)
			}
		}
		return
	}
	if chs {
		pos++
		res[pos] = u
		for _, v := range to[u] {
			if v != fa && deg[v] > 1 {
				dfs(v, u, false)
			}
		}
		for _, v := range to[u] {
			if v != fa && deg[v] == 1 {
				pos++
				res[pos] = v
			}
		}
	} else {
		for _, v := range to[u] {
			if v != fa && deg[v] == 1 {
				pos++
				res[pos] = v
			}
		}
		for _, v := range to[u] {
			if v != fa && deg[v] > 1 {
				dfs(v, u, true)
			}
		}
		pos++
		res[pos] = u
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	to = make([][]int, n+1)
	deg = make([]int, n+1)
	res = make([]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		to[u] = append(to[u], v)
		to[v] = append(to[v], u)
		deg[u]++
		deg[v]++
	}
	if n == 2 {
		fmt.Fprintln(out, "Yes")
		fmt.Fprintln(out, "1 2")
		out.Flush()
		return
	}
	for i := 1; i <= n; i++ {
		if deg[i] > 1 {
			rt = i
			break
		}
	}
	dfs(rt, 0, false)
	fmt.Fprintln(out, "Yes")
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, res[i])
	}
	out.WriteByte('\n')
	out.Flush()
}
