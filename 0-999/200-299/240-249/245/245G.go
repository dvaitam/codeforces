package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func check(n int, adj [][]int, N int) int {
	a := make([]int, N)
	// mark neighbors and self
	for _, nei := range adj[n] {
		a[nei] = 1
	}
	a[n] = 1
	m := 0
	// for each neighbor's neighbor, update counts
	for _, nei := range adj[n] {
		for _, nn := range adj[nei] {
			if a[nn]%2 == 0 {
				a[nn] += 2
				if a[nn] > m {
					m = a[nn]
				}
			}
		}
	}
	// count nodes with max value
	r := 0
	for i := 0; i < N; i++ {
		if a[i] == m {
			r++
		}
	}
	return r
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var K int
	if _, err := fmt.Fscan(in, &K); err != nil {
		return
	}
	m := make(map[string]int)
	adj := make([][]int, 0)
	N := 0
	for i := 0; i < K; i++ {
		var s, t string
		fmt.Fscan(in, &s, &t)
		id, ok := m[s]
		if !ok {
			id = N
			m[s] = N
			adj = append(adj, nil)
			N++
		}
		jd, ok2 := m[t]
		if !ok2 {
			jd = N
			m[t] = N
			adj = append(adj, nil)
			N++
		}
		adj[id] = append(adj[id], jd)
		adj[jd] = append(adj[jd], id)
	}
	// output
	fmt.Println(N)
	// sort keys lex
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		id := m[k]
		fmt.Printf("%s %d\n", k, check(id, adj, N))
	}
}
