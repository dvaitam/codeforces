package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	values := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &values[i])
	}

	degree := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		degree[u]++
		degree[v]++
	}

	// In a rooted tree with root 1, leaves have degree 1 except the root itself.
	leaves := 0
	for i := 2; i <= n; i++ {
		if degree[i] == 1 {
			leaves++
		}
	}

	// Naive strategy: change each leaf value individually.
	fmt.Fprintln(out, leaves)
}
