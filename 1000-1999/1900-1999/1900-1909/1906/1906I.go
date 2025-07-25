package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
	}

	// Hopcroft-Karp for maximum bipartite matching
	matchL := make([]int, n)
	matchR := make([]int, n)
	for i := 0; i < n; i++ {
		matchL[i] = -1
		matchR[i] = -1
	}
	dist := make([]int, n)

	bfs := func() bool {
		queue := make([]int, 0)
		for i := 0; i < n; i++ {
			if matchL[i] == -1 {
				dist[i] = 0
				queue = append(queue, i)
			} else {
				dist[i] = -1
			}
		}
		found := false
		for head := 0; head < len(queue); head++ {
			u := queue[head]
			for _, v := range adj[u] {
				m := matchR[v]
				if m != -1 && dist[m] == -1 {
					dist[m] = dist[u] + 1
					queue = append(queue, m)
				}
				if m == -1 {
					found = true
				}
			}
		}
		return found
	}

	var dfs func(int) bool
	dfs = func(u int) bool {
		for _, v := range adj[u] {
			m := matchR[v]
			if m == -1 || (dist[m] == dist[u]+1 && dfs(m)) {
				matchL[u] = v
				matchR[v] = u
				return true
			}
		}
		dist[u] = -1
		return false
	}

	matchCount := 0
	for bfs() {
		for i := 0; i < n; i++ {
			if matchL[i] == -1 && dfs(i) {
				matchCount++
			}
		}
	}

	// Construct path cover from matching
	next := make([]int, n)
	prev := make([]int, n)
	for i := 0; i < n; i++ {
		next[i] = -1
		prev[i] = -1
	}
	for u := 0; u < n; u++ {
		v := matchL[u]
		if v != -1 {
			next[u] = v
			prev[v] = u
		}
	}

	starts := make([]int, 0)
	tails := make([]int, 0)
	for i := 0; i < n; i++ {
		if prev[i] == -1 {
			starts = append(starts, i)
			cur := i
			for next[cur] != -1 {
				cur = next[cur]
			}
			tails = append(tails, cur)
		}
	}

	k := len(starts)
	fmt.Fprintln(writer, k-1)
	for i := 0; i < k-1; i++ {
		fmt.Fprintln(writer, tails[i]+1, starts[i+1]+1)
	}
}
