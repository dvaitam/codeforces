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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	adj := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		adj[i] = make([]bool, n+1)
	}

	for i := 1; i <= n; i++ {
		var cnt int
		fmt.Fscan(reader, &cnt)
		for j := 0; j < cnt; j++ {
			var x int
			fmt.Fscan(reader, &x)
			if x >= 1 && x <= n {
				adj[i][x] = true
			}
		}
	}

	type pair struct{ u, v int }
	var res []pair
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if adj[i][j] && !adj[j][i] {
				res = append(res, pair{j, i})
			}
		}
	}

	fmt.Fprintln(writer, len(res))
	for _, p := range res {
		fmt.Fprintf(writer, "%d %d\n", p.u, p.v)
	}
}
