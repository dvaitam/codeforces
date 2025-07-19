package main

import (
	"bufio"
	"fmt"
	"os"
)

var rd = bufio.NewReader(os.Stdin)
var wr = bufio.NewWriter(os.Stdout)

func readInt() int {
	c, err := rd.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = rd.ReadByte()
	}
	if err != nil {
		return 0
	}
	sign := 1
	if c == '-' {
		sign = -1
		c, _ = rd.ReadByte()
	}
	x := 0
	for err == nil && c >= '0' && c <= '9' {
		x = x*10 + int(c-'0')
		c, err = rd.ReadByte()
	}
	return x * sign
}

type edge struct {
	to int
	w  int64
}

func main() {
	defer wr.Flush()
	n := readInt()
	m := readInt()
	adj := make([][]edge, n+1)
	for i := 0; i < n-1; i++ {
		u := readInt()
		v := readInt()
		z := int64(readInt())
		adj[u] = append(adj[u], edge{v, z})
		adj[v] = append(adj[v], edge{u, z})
	}
	// find path from 1 to n using DFS
	parent := make([]int, n+1)
	wpar := make([]int64, n+1)
	stack := make([]int, 0, n)
	stack = append(stack, 1)
	parent[1] = 0
	found := false
	for len(stack) > 0 && !found {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if u == n {
			found = true
			break
		}
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			wpar[e.to] = e.w
			stack = append(stack, e.to)
		}
	}
	// build path and weights
	var path []int
	var wts []int64
	for u := n; u != 0; u = parent[u] {
		path = append(path, u)
		wts = append(wts, wpar[u])
	}
	L := len(path)
	nodes := make([]int, L)
	weight := make([]int64, L)
	for i := 0; i < L; i++ {
		nodes[i] = path[L-1-i]
		weight[i] = wts[L-1-i]
	}
	// prefix sum on path
	f := make([]int64, L)
	for i := 1; i < L; i++ {
		f[i] = f[i-1] + weight[i]
	}
	total := f[L-1]
	// h values
	h := make([]int64, L)
	for i := 0; i < L; i++ {
		h[i] = total - f[i]
	}
	// mark visited and compute branch weights g
	visited := make([]bool, n+1)
	for _, u := range nodes {
		visited[u] = true
	}
	g := make([]int64, L)
	for i, u := range nodes {
		for _, e := range adj[u] {
			if !visited[e.to] {
				visited[e.to] = true
				g[i] = e.w
				break
			}
		}
	}
	// initial ans
	ans := int64(2000000001)
	for i := 1; i <= n; i++ {
		if !visited[i] {
			ans = 0
			break
		}
	}
	// compute u array
	uarr := make([]int64, L)
	for i := 0; i < L; i++ {
		if i > 0 {
			uarr[i] = uarr[i-1]
		}
		if g[i] > 0 {
			val := f[i] + g[i]
			if val > uarr[i] {
				uarr[i] = val
			}
		}
	}
	// evaluate ans
	for i := L - 1; i >= 1; i-- {
		if g[i] > 0 && uarr[i-1] > 0 {
			tmp := total - uarr[i-1] - h[i] - g[i]
			if tmp < ans {
				ans = tmp
			}
		}
	}
	for i := 0; i+2 < L; i++ {
		tmp := total - f[i] - h[i+2]
		if tmp < ans {
			ans = tmp
		}
	}
	for i := 0; i+1 < L; i++ {
		if g[i] > 0 {
			tmp := total - f[i] - h[i+1] - g[i]
			if tmp < ans {
				ans = tmp
			}
		}
		if g[i+1] > 0 {
			tmp := total - f[i] - h[i+1] - g[i+1]
			if tmp < ans {
				ans = tmp
			}
		}
	}
	if ans < 0 {
		ans = 0
	}
	// answer queries
	for i := 0; i < m; i++ {
		x := int64(readInt())
		if x >= ans {
			fmt.Fprintln(wr, total)
		} else {
			fmt.Fprintln(wr, total-ans+x)
		}
	}
}
