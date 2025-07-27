package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		a--
		b--
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	topics := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &topics[i])
	}

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		if topics[order[i]] != topics[order[j]] {
			return topics[order[i]] < topics[order[j]]
		}
		return order[i] < order[j]
	})

	seen := make([]int, n+2)
	curIter := 1

	for _, v := range order {
		tV := topics[v]
		cnt := 0
		for _, nb := range adj[v] {
			tN := topics[nb]
			if tN == tV {
				fmt.Fprintln(writer, -1)
				return
			}
			if tN < tV {
				if seen[tN] != curIter {
					seen[tN] = curIter
					cnt++
				}
			}
		}
		if cnt != tV-1 {
			fmt.Fprintln(writer, -1)
			return
		}
		curIter++
	}

	for i, v := range order {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v+1)
	}
	writer.WriteByte('\n')
}
