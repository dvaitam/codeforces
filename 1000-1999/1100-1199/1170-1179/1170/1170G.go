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

	var n, m int
	fmt.Fscan(reader, &n, &m)
	adj := make([][]edge, n+1)
	deg := make([]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], edge{to: v, id: i})
		adj[v] = append(adj[v], edge{to: u, id: i})
		deg[u]++
		deg[v]++
	}
	for v := 1; v <= n; v++ {
		if deg[v]%2 != 0 {
			fmt.Fprintln(writer, "NO")
			return
		}
	}
	// Hierholzer per component
	ptr := make([]int, n+1)
	usedEdge := make([]bool, m)
	var cycles [][]int
	// temp buffers
	stack := make([]int, 0, m+1)
	path := make([]int, 0, m+1)
	// prepare id mapping for cycle decomposition
	id := make([]int, n+1)
	for i := range id {
		id[i] = -1
	}
	for v0 := 1; v0 <= n; v0++ {
		if ptr[v0] < len(adj[v0]) {
			// start Euler tour at v0
			stack = stack[:0]
			path = path[:0]
			stack = append(stack, v0)
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				// find next unused edge
				for ptr[u] < len(adj[u]) && usedEdge[adj[u][ptr[u]].id] {
					ptr[u]++
				}
				if ptr[u] == len(adj[u]) {
					// dead end
					path = append(path, u)
					stack = stack[:len(stack)-1]
				} else {
					e := adj[u][ptr[u]]
					usedEdge[e.id] = true
					// advance pointer
					ptr[u]++
					stack = append(stack, e.to)
				}
			}
			// path is reverse of tour
			// reverse to get P
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
           // decompose path into simple cycles
           // id tracks position in stack for current active vertices
           // touched records vertices that were assigned id during this component
           touched := make([]int, 0, len(path))
           st := make([]int, 0, len(path))
           for _, vv := range path {
               if id[vv] == -1 {
                   st = append(st, vv)
                   id[vv] = len(st) - 1
                   touched = append(touched, vv)
               } else {
                   idx := id[vv]
                   // form cycle st[idx..] + vv
                   cyc := make([]int, 0, len(st)-idx+1)
                   for k := idx; k < len(st); k++ {
                       cyc = append(cyc, st[k])
                   }
                   cyc = append(cyc, vv)
                   cycles = append(cycles, cyc)
                   // reset id for removed vertices
                   for k := idx + 1; k < len(st); k++ {
                       id[st[k]] = -1
                   }
                   // keep st[0..idx]
                   st = st[:idx+1]
               }
           }
           // clear id for this component
           for _, v := range touched {
               id[v] = -1
           }
		}
	}
	// output
	fmt.Fprintln(writer, "YES")
	fmt.Fprintln(writer, len(cycles))
	for _, cyc := range cycles {
		fmt.Fprint(writer, len(cyc))
		for _, x := range cyc {
			fmt.Fprint(writer, " ", x)
		}
		fmt.Fprintln(writer)
	}
}

type edge struct {
	to int
	id int
}
