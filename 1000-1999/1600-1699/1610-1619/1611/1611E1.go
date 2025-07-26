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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		friends := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &friends[i])
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		const INF = int(1e9)
		distFriend := make([]int, n+1)
		for i := range distFriend {
			distFriend[i] = INF
		}
		queue := make([]int, 0, n)
		for _, x := range friends {
			distFriend[x] = 0
			queue = append(queue, x)
		}
		head := 0
		for head < len(queue) {
			u := queue[head]
			head++
			for _, v := range adj[u] {
				if distFriend[v] == INF {
					distFriend[v] = distFriend[u] + 1
					queue = append(queue, v)
				}
			}
		}

		distVlad := make([]int, n+1)
		for i := range distVlad {
			distVlad[i] = INF
		}
		queue = queue[:0]
		distVlad[1] = 0
		queue = append(queue, 1)
		head = 0
		escaped := false
		for head < len(queue) {
			u := queue[head]
			head++
			if u != 1 && len(adj[u]) == 1 {
				escaped = true
				break
			}
			for _, v := range adj[u] {
				if distVlad[v] == INF && distVlad[u]+1 < distFriend[v] {
					distVlad[v] = distVlad[u] + 1
					queue = append(queue, v)
				}
			}
		}
		if escaped {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
