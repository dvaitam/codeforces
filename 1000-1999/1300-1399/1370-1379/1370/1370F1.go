package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(w *bufio.Writer, r *bufio.Reader, nodes []int) (int, int) {
	fmt.Fprint(w, "? ")
	fmt.Fprint(w, len(nodes))
	for _, v := range nodes {
		fmt.Fprint(w, " ", v)
	}
	fmt.Fprintln(w)
	w.Flush()
	var node, dist int
	fmt.Fscan(r, &node, &dist)
	return node, dist
}

func bfs(start int, g [][]int) ([]int, []int) {
	n := len(g) - 1
	dist := make([]int, n+1)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				parent[to] = v
				q = append(q, to)
			}
		}
	}
	return dist, parent
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		all := make([]int, n)
		for i := 1; i <= n; i++ {
			all[i-1] = i
		}
		root, distSF := query(writer, reader, all)
		dRoot, _ := bfs(root, g)

		var cand []int
		for i := 1; i <= n; i++ {
			if dRoot[i] == distSF {
				cand = append(cand, i)
			}
		}
		node1, _ := query(writer, reader, cand)
		dNode1, _ := bfs(node1, g)
		var cand2 []int
		for i := 1; i <= n; i++ {
			if dNode1[i] == distSF {
				cand2 = append(cand2, i)
			}
		}
		node2, _ := query(writer, reader, cand2)
		fmt.Fprintf(writer, "! %d %d\n", node1, node2)
		writer.Flush()
		var verdict string
		fmt.Fscan(reader, &verdict)
		if verdict != "Correct" && verdict != "" {
			return
		}
	}
}
