package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ u, v int }

func readInt(r *bufio.Reader) int {
	sign := 1
	c, _ := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = r.ReadByte()
	}
	x := int(c - '0')
	for {
		c, err := r.ReadByte()
		if err != nil || c < '0' || c > '9' {
			break
		}
		x = x*10 + int(c-'0')
	}
	return sign * x
}

func addEdge(adj map[int]map[int]bool, u, v int) {
	if adj[u] == nil {
		adj[u] = make(map[int]bool)
	}
	if adj[v] == nil {
		adj[v] = make(map[int]bool)
	}
	adj[u][v] = true
	adj[v][u] = true
}

func removeEdge(adj map[int]map[int]bool, u, v int) {
	if nb, ok := adj[u]; ok {
		delete(nb, v)
		if len(nb) == 0 {
			delete(adj, u)
		}
	}
	if nb, ok := adj[v]; ok {
		delete(nb, u)
		if len(nb) == 0 {
			delete(adj, v)
		}
	}
}

func isBipartite(adj map[int]map[int]bool, n int) bool {
	color := make([]int, n+1)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if color[i] == 0 && len(adj[i]) > 0 {
			color[i] = 1
			q = append(q, i)
			for len(q) > 0 {
				x := q[0]
				q = q[1:]
				for y := range adj[x] {
					if color[y] == 0 {
						color[y] = 3 - color[x]
						q = append(q, y)
					} else if color[y] == color[x] {
						return false
					}
				}
			}
		}
	}
	return true
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	n := readInt(r)
	m := readInt(r)
	k := readInt(r)
	qn := readInt(r)

	edges := make([]Edge, m+1)
	for i := 1; i <= m; i++ {
		u := readInt(r)
		v := readInt(r)
		edges[i] = Edge{u, v}
	}

	curColor := make([]int, m+1)
	graphs := make([]map[int]map[int]bool, k+1)
	for i := 1; i <= k; i++ {
		graphs[i] = make(map[int]map[int]bool)
	}

	for qi := 0; qi < qn; qi++ {
		e := readInt(r)
		c := readInt(r)
		old := curColor[e]
		if old == c {
			fmt.Fprintln(w, "YES")
			continue
		}
		u, v := edges[e].u, edges[e].v
		if old != 0 {
			removeEdge(graphs[old], u, v)
		}
		addEdge(graphs[c], u, v)
		if isBipartite(graphs[c], n) {
			curColor[e] = c
			fmt.Fprintln(w, "YES")
		} else {
			removeEdge(graphs[c], u, v)
			if old != 0 {
				addEdge(graphs[old], u, v)
			}
			fmt.Fprintln(w, "NO")
		}
	}
}
