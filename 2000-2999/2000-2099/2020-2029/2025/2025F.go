package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	id int
}

func assignEdge(id, vertex int, xs, ys []int, ans []byte, used []bool) {
	used[id] = true
	if xs[id] == vertex {
		ans[id] = 'x'
	} else if ys[id] == vertex {
		ans[id] = 'y'
	} else {
		ans[id] = 'x'
	}
}

func orientComponent(root int, adj [][]edge, need []uint8, parent []int, parentEdge []int, visited []bool, edgeUsed []bool, ans []byte, xs, ys []int) {
	type stackItem struct {
		node int
		idx  int
	}
	stack := []stackItem{{root, 0}}
	visited[root] = true
	parent[root] = -1
	parentEdge[root] = -1
	order := make([]int, 0)

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.node
		if top.idx < len(adj[u]) {
			e := adj[u][top.idx]
			top.idx++
			v := e.to
			id := e.id
			if parent[u] == v && parentEdge[u] == id {
				continue
			}
			if !visited[v] {
				visited[v] = true
				parent[v] = u
				parentEdge[v] = id
				stack = append(stack, stackItem{v, 0})
			} else if !edgeUsed[id] {
				assignEdge(id, u, xs, ys, ans, edgeUsed)
				need[u] ^= 1
			}
		} else {
			order = append(order, u)
			stack = stack[:len(stack)-1]
		}
	}

	for _, u := range order {
		if u == root {
			continue
		}
		id := parentEdge[u]
		if edgeUsed[id] {
			continue
		}
		if need[u] == 1 {
			assignEdge(id, u, xs, ys, ans, edgeUsed)
			need[u] ^= 1
		} else {
			assignEdge(id, parent[u], xs, ys, ans, edgeUsed)
			need[parent[u]] ^= 1
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	adj := make([][]edge, n)
	xs := make([]int, q)
	ys := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
		xs[i]--
		ys[i]--
		adj[xs[i]] = append(adj[xs[i]], edge{ys[i], i})
		adj[ys[i]] = append(adj[ys[i]], edge{xs[i], i})
	}

	visitedComp := make([]bool, n)
	visited := make([]bool, n)
	need := make([]uint8, n)
	parent := make([]int, n)
	parentEdge := make([]int, n)
	edgeUsed := make([]bool, q)
	ans := make([]byte, q)

	queue := make([]int, 0)

	for start := 0; start < n; start++ {
		if visitedComp[start] || len(adj[start]) == 0 {
			continue
		}
		queue = queue[:0]
		queue = append(queue, start)
		visitedComp[start] = true
		compNodes := make([]int, 0)
		totalDeg := 0
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			compNodes = append(compNodes, u)
			totalDeg += len(adj[u])
			for _, e := range adj[u] {
				v := e.to
				if !visitedComp[v] {
					visitedComp[v] = true
					queue = append(queue, v)
				}
			}
		}
		for _, u := range compNodes {
			need[u] = 0
		}
		if (totalDeg/2)%2 == 1 {
			need[start] = 1
		}
		if !visited[start] {
			orientComponent(start, adj, need, parent, parentEdge, visited, edgeUsed, ans, xs, ys)
		}
	}

	for i := 0; i < q; i++ {
		if !edgeUsed[i] {
			assignEdge(i, xs[i], xs, ys, ans, edgeUsed)
		}
	}

	val := make([]int, n)
	for i := 0; i < q; i++ {
		var vertex int
		if ans[i] == 'y' {
			vertex = ys[i]
		} else {
			vertex = xs[i]
		}
		if val[vertex] > 0 {
			val[vertex]--
			fmt.Fprintf(out, "%c-\n", ans[i])
		} else {
			val[vertex]++
			fmt.Fprintf(out, "%c+\n", ans[i])
		}
	}
}

