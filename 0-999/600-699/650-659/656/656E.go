package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &dist[i][j])
		}
	}

	// Floyd-Warshall algorithm
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	maxDist := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if dist[i][j] > maxDist {
				maxDist = dist[i][j]
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, maxDist)
	writer.Flush()
}
