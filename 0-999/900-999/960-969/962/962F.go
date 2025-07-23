package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	id int
}

var (
	n, m     int
	g        [][]Edge
	disc     []int
	low      []int
	timer    int
	isBridge []bool
)

func dfs(u, parentEdge int) {
	timer++
	disc[u] = timer
	low[u] = timer
	for _, e := range g[u] {
		if e.id == parentEdge {
			continue
		}
		v := e.to
		if disc[v] == 0 {
			dfs(v, e.id)
			if low[v] < low[u] {
				low[u] = low[v]
			}
			if low[v] > disc[u] {
				isBridge[e.id] = true
			}
		} else {
			if disc[v] < low[u] {
				low[u] = disc[v]
			}
		}
	}
}

var (
	visited     []bool
	edgeVisited []bool
	ans         []int
)

func dfs2(u int, vertexCount *int, edges *[]int) {
	visited[u] = true
	*vertexCount = *vertexCount + 1
	for _, e := range g[u] {
		if isBridge[e.id] {
			continue
		}
		if !edgeVisited[e.id] {
			edgeVisited[e.id] = true
			*edges = append(*edges, e.id)
		}
		if !visited[e.to] {
			dfs2(e.to, vertexCount, edges)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	for i := 1; i <= m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], Edge{v, i})
		g[v] = append(g[v], Edge{u, i})
	}
	disc = make([]int, n+1)
	low = make([]int, n+1)
	isBridge = make([]bool, m+1)
	for i := 1; i <= n; i++ {
		if disc[i] == 0 {
			dfs(i, -1)
		}
	}

	visited = make([]bool, n+1)
	edgeVisited = make([]bool, m+1)
	for i := 1; i <= n; i++ {
		if !visited[i] {
			tempEdges := []int{}
			count := 0
			dfs2(i, &count, &tempEdges)
			if len(tempEdges) == count {
				ans = append(ans, tempEdges...)
			}
		}
	}
	sort.Ints(ans)
	fmt.Fprintln(writer, len(ans))
	for i, id := range ans {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, id)
	}
	if len(ans) > 0 {
		fmt.Fprintln(writer)
	}
}
