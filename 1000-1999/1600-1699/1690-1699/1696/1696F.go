package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct{ l, r int }

func getGraphByEdge(e edge, d [][][]bool, visited []bool, edges *[]edge, paths [][]int) bool {
	n := len(d)
	// initialize path for this edge
	paths[e.l][e.r] = 1
	paths[e.r][e.l] = 1
	// queue for BFS
	q := []edge{e}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		l, r := cur.l, cur.r
		for i := 0; i < n; i++ {
			if d[l][r][i] {
				if visited[i] && i != l {
					return false
				}
				// add new edge
				q = append(q, edge{r, i})
				visited[i] = true
				*edges = append(*edges, edge{r, i})
				// update paths
				for j := 0; j < n; j++ {
					if paths[j][r] >= 0 && j != i {
						paths[j][i] = paths[j][r] + 1
						paths[i][j] = paths[j][r] + 1
					}
				}
			}
		}
	}
	return true
}

func checkTree(d [][][]bool, paths [][]int) bool {
	n := len(d)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := i + 1; k < n; k++ {
				if d[i][j][k] {
					if paths[i][j] != paths[j][k] {
						return false
					}
				} else {
					if paths[i][j] == paths[j][k] {
						return false
					}
				}
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for tc := 0; tc < T; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		// d[i][j][k]
		d := make([][][]bool, n)
		for i := 0; i < n; i++ {
			d[i] = make([][]bool, n)
			for j := 0; j < n; j++ {
				d[i][j] = make([]bool, n)
			}
		}
		// read input strings
		for i := 0; i < n-1; i++ {
			for j := i; j < n-1; j++ {
				var s string
				fmt.Fscan(reader, &s)
				for k := 0; k < len(s); k++ {
					if s[k] == '1' {
						d[i][k][j+1] = true
						d[j+1][k][i] = true
					}
				}
			}
		}
		found := false
		// try initial edge from 0 to i
		for i := 1; i < n && !found; i++ {
			visited := make([]bool, n)
			visited[0], visited[i] = true, true
			edges := make([]edge, 0, n-1)
			// init paths
			paths := make([][]int, n)
			for u := 0; u < n; u++ {
				paths[u] = make([]int, n)
				for v := 0; v < n; v++ {
					paths[u][v] = -1
				}
				paths[u][u] = 0
			}
			// first BFS
			if getGraphByEdge(edge{0, i}, d, visited, &edges, paths) {
				// second BFS
				if getGraphByEdge(edge{i, 0}, d, visited, &edges, paths) {
					// add initial edge
					edges = append(edges, edge{0, i})
					if len(edges) == n-1 && checkTree(d, paths) {
						fmt.Fprintln(writer, "Yes")
						for _, e := range edges {
							// output 1-based
							fmt.Fprintln(writer, e.l+1, e.r+1)
						}
						found = true
					}
				}
			}
		}
		if !found {
			fmt.Fprintln(writer, "No")
		}
	}
}
