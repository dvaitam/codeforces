package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

const maxLog = 20

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]int, n)
	parent := make([][maxLog]int, n)
	depth := make([]int, n)
	for i := 1; i < n; i++ {
		var p int
		fmt.Fscan(in, &p)
		p--
		g[p] = append(g[p], i)
		g[i] = append(g[i], p)
		parent[i][0] = p
	}
	// dfs to compute depth and parents
	stack := []int{0}
	order := []int{0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if to == parent[v][0] && v != 0 {
				continue
			}
			parent[to][0] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
			order = append(order, to)
		}
	}
	for j := 1; j < maxLog; j++ {
		for i := 0; i < n; i++ {
			parent[i][j] = parent[parent[i][j-1]][j-1]
		}
	}
	// lca function
	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := 0; i < maxLog; i++ {
			if diff>>i&1 == 1 {
				a = parent[a][i]
			}
		}
		if a == b {
			return a
		}
		for i := maxLog - 1; i >= 0; i-- {
			if parent[a][i] != parent[b][i] {
				a = parent[a][i]
				b = parent[b][i]
			}
		}
		return parent[a][0]
	}

	var m int
	fmt.Fscan(in, &m)
	routeNodes := make([][]int, m)
	cityRoutes := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		gNode := lca(a, b)
		path := []int{}
		x := a
		for x != gNode {
			path = append(path, x)
			x = parent[x][0]
		}
		tmp := []int{}
		y := b
		for y != gNode {
			tmp = append(tmp, y)
			y = parent[y][0]
		}
		path = append(path, gNode)
		for i := len(tmp) - 1; i >= 0; i-- {
			path = append(path, tmp[i])
		}
		routeNodes[i] = path
		for _, v := range path {
			cityRoutes[v] = append(cityRoutes[v], i)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		// 0-1 BFS
		tot := n + m
		dist := make([]int, tot)
		for i := range dist {
			dist[i] = -1
		}
		dq := list.New()
		dist[u] = 0
		dq.PushBack(u)
		for dq.Len() > 0 {
			e := dq.Front()
			dq.Remove(e)
			x := e.Value.(int)
			if x == v {
				break
			}
			if x < n {
				for _, r := range cityRoutes[x] {
					y := n + r
					if dist[y] == -1 || dist[y] > dist[x]+1 {
						dist[y] = dist[x] + 1
						dq.PushBack(y)
					}
				}
			} else {
				r := x - n
				for _, y := range routeNodes[r] {
					if dist[y] == -1 || dist[y] > dist[x] {
						dist[y] = dist[x]
						dq.PushFront(y)
					}
				}
			}
		}
		fmt.Fprintln(out, dist[v])
	}
}
