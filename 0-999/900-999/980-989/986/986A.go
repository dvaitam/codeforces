package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k, s int
	if _, err := fmt.Fscan(in, &n, &m, &k, &s); err != nil {
		return
	}
	types := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &types[i])
		types[i]--
	}
	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	dist := make([][]int, k)
	q := make([]int, 0, n)
	for t := 0; t < k; t++ {
		dist[t] = make([]int, n)
		for i := 0; i < n; i++ {
			dist[t][i] = -1
		}
		q = q[:0]
		for i := 0; i < n; i++ {
			if types[i] == t {
				dist[t][i] = 0
				q = append(q, i)
			}
		}
		for head := 0; head < len(q); head++ {
			v := q[head]
			nd := dist[t][v] + 1
			for _, to := range adj[v] {
				if dist[t][to] == -1 {
					dist[t][to] = nd
					q = append(q, to)
				}
			}
		}
	}

	tmp := make([]int, k)
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		for t := 0; t < k; t++ {
			tmp[t] = dist[t][i]
		}
		sort.Ints(tmp)
		sum := 0
		for j := 0; j < s; j++ {
			sum += tmp[j]
		}
		ans[i] = sum
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
