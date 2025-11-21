package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	idx int
}

type Added struct {
	u, v int
}

func hopcroftKarp(n int, adj [][]Edge, colored []bool) (matchL []int, edgeL []int) {
	matchL = make([]int, n+1)
	matchR := make([]int, n+1)
	edgeL = make([]int, n+1)
	edgeR := make([]int, n+1)
	dist := make([]int, n+1)
	bfs := func() bool {
		q := make([]int, 0)
		for u := 1; u <= n; u++ {
			if matchL[u] == 0 {
				dist[u] = 0
				q = append(q, u)
			} else {
				dist[u] = -1
			}
		}
		found := false
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, e := range adj[u] {
				if colored[e.idx] {
					continue
				}
				v := e.to
				if matchR[v] != 0 {
					if dist[matchR[v]] == -1 {
						dist[matchR[v]] = dist[u] + 1
						q = append(q, matchR[v])
					}
				} else {
					found = true
				}
			}
		}
		return found
	}
	var dfs func(u int) bool
	dfs = func(u int) bool {
		for _, e := range adj[u] {
			if colored[e.idx] {
				continue
			}
			v := e.to
			if matchR[v] == 0 || (dist[matchR[v]] == dist[u]+1 && dfs(matchR[v])) {
				matchL[u] = v
				matchR[v] = u
				edgeL[u] = e.idx
				edgeR[v] = e.idx
				return true
			}
		}
		dist[u] = -1
		return false
	}
	for bfs() {
		for u := 1; u <= n; u++ {
			if matchL[u] == 0 {
				dfs(u)
			}
		}
	}
	return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edgesU := make([]int, m+1)
	edgesV := make([]int, m+1)
	adj := make([][]Edge, n+1)
	outDeg := make([]int, n+1)
	inDeg := make([]int, n+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &edgesU[i], &edgesV[i])
		adj[edgesU[i]] = append(adj[edgesU[i]], Edge{edgesV[i], i})
		outDeg[edgesU[i]]++
		inDeg[edgesV[i]]++
	}

	maxDeg := 0
	for i := 1; i <= n; i++ {
		if outDeg[i] > maxDeg {
			maxDeg = outDeg[i]
		}
		if inDeg[i] > maxDeg {
			maxDeg = inDeg[i]
		}
	}
	if maxDeg == 0 {
		maxDeg = 1
	}

	colors := make([]int, m+1)
	colored := make([]bool, m+1)
	for c := 1; c <= maxDeg; c++ {
		matchL, edgeL := hopcroftKarp(n, adj, colored)
		matched := 0
		for u := 1; u <= n; u++ {
			if matchL[u] != 0 {
				idx := edgeL[u]
				colors[idx] = c
				colored[idx] = true
				matched++
			}
		}
		if matched == 0 {
			for u := 1; u <= n; u++ {
				if matchL[u] == 0 {
					// ensure progress by adding dummy edge (not needed)
				}
			}
		}
	}

	added := make([]Added, 0)
	totalEdges := m
	for c := 1; c <= maxDeg; c++ {
		next := make([]int, n+1)
		prev := make([]int, n+1)
		for idx := 1; idx <= totalEdges; idx++ {
			if colors[idx] == c {
				u := edgesU[idx]
				v := edgesV[idx]
				next[u] = v
				prev[v] = u
			}
		}
		visited := make([]bool, n+1)
		heads := make([]int, 0)
		tails := make([]int, 0)
		for i := 1; i <= n; i++ {
			if prev[i] == 0 && !visited[i] {
				cur := i
				for !visited[cur] {
					visited[cur] = true
					if next[cur] == 0 {
						break
					}
					cur = next[cur]
				}
				heads = append(heads, i)
				tails = append(tails, cur)
			}
		}
		for i := 1; i <= n; i++ {
			if !visited[i] {
				cur := i
				for prev[cur] != 0 && !visited[prev[cur]] {
					cur = prev[cur]
				}
				head := cur
				for !visited[cur] {
					visited[cur] = true
					if next[cur] == 0 {
						break
					}
					cur = next[cur]
				}
				heads = append(heads, head)
				tails = append(tails, cur)
			}
		}
		k := len(heads)
		if k == 0 {
			heads = append(heads, 1)
			tails = append(tails, 1)
			k = 1
		}
		for i := 0; i < k; i++ {
			from := tails[i]
			to := heads[(i+1)%k]
			if next[from] == 0 {
				totalEdges++
				edgesU = append(edgesU, from)
				edgesV = append(edgesV, to)
				colors = append(colors, c)
				added = append(added, Added{from, to})
				next[from] = to
				prev[to] = from
			}
		}
	}

	fmt.Fprintln(out, len(added))
	for _, e := range added {
		fmt.Fprintln(out, e.u, e.v)
	}
	fmt.Fprintln(out, maxDeg)
	for i := 1; i <= totalEdges; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, colors[i])
	}
	fmt.Fprintln(out)
}
