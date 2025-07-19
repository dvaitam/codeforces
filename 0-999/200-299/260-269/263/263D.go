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

	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// build a simple path by greedy DFS until no unvisited neighbor
	visited := make([]bool, n+1)
	c := make([]int, n+2)
	cnt := 1
	cur := 1
	visited[1] = true
	c[1] = 1
	for {
		found := false
		for _, nb := range adj[cur] {
			if !visited[nb] {
				visited[nb] = true
				cnt++
				c[cnt] = nb
				cur = nb
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	// find a cycle by connecting end back to earlier node
	end := c[cnt]
	for i := 1; i <= cnt; i++ {
		for _, nb := range adj[c[i]] {
			if nb == end {
				length := cnt - i + 1
				fmt.Fprintln(writer, length)
				for j := i; j <= cnt; j++ {
					if j > i {
						writer.WriteByte(' ')
					}
					fmt.Fprint(writer, c[j])
				}
				fmt.Fprintln(writer)
				return
			}
		}
	}
}
