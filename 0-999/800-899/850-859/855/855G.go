package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct{ u, v int }

func bfs(graph [][]int, start, target int, block map[Pair]bool) ([]Pair, bool) {
	n := len(graph) - 1
	visited := make([]bool, n+1)
	parent := make([]int, n+1)
	q := []int{start}
	visited[start] = true
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur == target {
			break
		}
		for _, nxt := range graph[cur] {
			if block[Pair{cur, nxt}] || block[Pair{nxt, cur}] {
				continue
			}
			if !visited[nxt] {
				visited[nxt] = true
				parent[nxt] = cur
				q = append(q, nxt)
			}
		}
	}
	if !visited[target] {
		return nil, false
	}
	var path []Pair
	x := target
	for x != start {
		p := parent[x]
		path = append(path, Pair{p, x})
		x = p
	}
	return path, true
}

func edgeDisjoint(graph [][]int, u, v, w int) bool {
	block := make(map[Pair]bool)
	path1, ok := bfs(graph, u, w, block)
	if !ok {
		return false
	}
	for _, e := range path1 {
		block[e] = true
		block[Pair{e.v, e.u}] = true
	}
	_, ok = bfs(graph, v, w, block)
	return ok
}

func countTriplets(graph [][]int) int64 {
	n := len(graph) - 1
	var ans int64
	for w := 1; w <= n; w++ {
		for u := 1; u <= n; u++ {
			if u == w {
				continue
			}
			for v := 1; v <= n; v++ {
				if v == w || v == u {
					continue
				}
				if edgeDisjoint(graph, u, v, w) {
					ans++
				}
			}
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	graph := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	var q int
	fmt.Fscan(in, &q)
	edges := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
	}

	ans := countTriplets(graph)
	fmt.Println(ans)
	for i := 0; i < q; i++ {
		e := edges[i]
		graph[e[0]] = append(graph[e[0]], e[1])
		graph[e[1]] = append(graph[e[1]], e[0])
		ans = countTriplets(graph)
		fmt.Println(ans)
	}
}
